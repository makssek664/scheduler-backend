package main;

import (
	"net/http"
	"time"
	"gorm.io/gorm"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
	"fmt"
	"errors"
	"strconv"

)

type HttpServer struct {
	s *http.Server;
	mux *mux.Router
	db *gorm.DB;
};

var DoesNotExistError = errors.New("httpServer: requested element does NOT exist")

type okS struct {
	Ok bool `json:"ok"`
}


func httpOkHandler(w http.ResponseWriter, r *http.Request) {
	ok := &okS{Ok : true};
	okb, _ := json.Marshal(ok)

	w.Write(okb)
}

func (s *HttpServer) authHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
        return
    }

	u := &User{}
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024);
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&u)

	if(err != nil) {
		log.Printf("JSON Decoding error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var count int64

	err = s.db.Session(&gorm.Session{}).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&User{}).Unscoped().Where("name = ?", u.Name).Count(&count)
		if res.Error != nil {
			return res.Error
		}

		if count == 0 {
			res = tx.Create(&u)
			if res != nil {
				return res.Error
			}
			res.Commit()
		} else {
			res = tx.Where("name = ?", u.Name).First(&u)
			if res != nil {
				return res.Error;
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("GORM TRANSACTION FAILURE: %v", err)
		http.Error(w, "Database op failed", http.StatusInternalServerError);
		return
	}
	rets, err := json.Marshal(u);
	if(err != nil) {
		log.Printf("JSON Encoding error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	w.Write(rets)
}

func (s *HttpServer) addEventToUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
    if r.Method == http.MethodOptions {
        return
    }

	e := &Event{}
	vars := mux.Vars(r)
	strid, ok := vars["user_id"]
	if !ok {
		http.Error(w, "Missing User ID param", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(strid, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}	

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024);
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&e)
	
	if err != nil {
		log.Printf("JSON Decoding error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest);
		return
	}
	if (e.Date == time.Time{}) {
		log.Printf("Empty date")
		http.Error(w, "Empty date in payload", http.StatusBadRequest)
	}

	var count int64
	err = s.db.Session(&gorm.Session{}).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&User{}).Where("id = ?", id).Count(&count)
		if res.Error != nil {
			return res.Error
		}
		

		if count == 0 {
			return DoesNotExistError	
		}

		e.UserID = uint(id)

		tx.Model(&Event{}).Create(&e)		

		if err != nil {
			return res.Error;
		}

		return nil
	})

	if err != nil {
		log.Printf("GORM Transaction Failure: %v", err);
		http.Error(w, "GORM Transaction failure", http.StatusBadRequest)
		return 
	}
	
	rets, err := json.Marshal(e)
	w.Write(rets)
}

func (s *HttpServer) getAllEventsUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
    if r.Method == http.MethodOptions {
        return
    }

	es := []Event{}	
	
	vars := mux.Vars(r)
	id, ok := vars["user_id"]
	if !ok {
		http.Error(w, "Missing User ID param", http.StatusBadRequest)
		return
	}

	var count int64
	err := s.db.Session(&gorm.Session{}).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&User{}).Where("id = ?", id).Count(&count)
		if res.Error != nil {
			return res.Error
		}
		if count == 0 {
			return DoesNotExistError
		}
		
		res = tx.Model(&Event{}).Where("user_id = ?", id).Find(&es)
		if res.Error != nil {
			return res.Error;
		}

		return nil
	})
	if err != nil {
		log.Printf("GORM Transaction Failure: %v", err);
		http.Error(w, "GORM Transaction failure", http.StatusBadRequest)
		return
	}
	
	rets, err := json.Marshal(es)
	if err != nil {
		log.Printf("JSON Marshaling error: %v", err);
		http.Error(w, "JSON Marshaling error", http.StatusInternalServerError);
	}
	w.Write(rets)
}

func (s *HttpServer) getEventByIDUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    if r.Method == http.MethodOptions {
        return
    }
	e := &Event{}	
	
	vars := mux.Vars(r)
	user_id, ok := vars["user_id"]
	if !ok {
		http.Error(w, "Missing User ID param", http.StatusBadRequest)
		return
	}
	event_id, ok := vars["event_id"]
	if !ok {
		http.Error(w, "Missing Event ID param", http.StatusBadRequest)
		return
	}


	var count int64
	err := s.db.Session(&gorm.Session{}).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&User{}).Where("id = ?", user_id).Count(&count)
		if res.Error != nil {
			return res.Error
		}
		if count == 0 {
			return DoesNotExistError
		}
		
		res = tx.Model(&Event{}).Where("user_id = ?", user_id).Where("id = ?", event_id).First(&e)
		if res.Error != nil {
			return res.Error;
		}

		return nil
	})
	if err != nil {
		log.Printf("GORM Transaction Failure: %v", err);
		http.Error(w, "GORM Transaction failure", http.StatusBadRequest)
		return
	}
	
	rets, err := json.Marshal(e)	
	if err != nil {
		log.Printf("JSON Marshaling error: %v", err);
		http.Error(w, "JSON Marshaling error", http.StatusInternalServerError);
	}
	w.Write(rets)
}

func (s *HttpServer) getEventByDateReqAndIDUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    if r.Method == http.MethodOptions {
        return
    }
	
	es := []Event{}
	vars := mux.Vars(r)
	
	user_id, ok := vars["user_id"]
	if !ok {
		http.Error(w, "Missing User ID param", http.StatusBadRequest)
		return
	}
	from, ok := vars["from"]
	if !ok {
		http.Error(w, "Missing FROM param", http.StatusBadRequest)
		return
	}
	until, ok := vars["until"]
	if !ok {
		http.Error(w, "Missing UNTIL param", http.StatusBadRequest)
		return
	}

	fromd, err := time.Parse("UTC", from);
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest);
		return
	}
	untild, err := time.Parse("UTC", until);
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest);
		return
	}

	var count int64
	err = s.db.Session(&gorm.Session{}).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&User{}).Where("id = ?", user_id).Count(&count)
		if res.Error != nil {
			return res.Error
		}
		if count == 0 {
			return DoesNotExistError
		}
		
		res = tx.Model(&Event{}).Where("user_id = ?", user_id).
		Where("date > ?", fromd).
		Where("date < ?", untild).
		Find(&es)

		if res.Error != nil {
			return res.Error;
		}

		return nil
	})
	if err != nil {
		log.Printf("GORM Transaction Failure: %v", err);
		http.Error(w, "GORM Transaction failure", http.StatusBadRequest)
		return
	}
	
	rets, err := json.Marshal(es)
	if err != nil {
		log.Printf("JSON Marshaling error: %v", err);
		http.Error(w, "JSON Marshaling error", http.StatusInternalServerError);
	}
	w.Write(rets)

}

func (s *HttpServer) setEventByIDUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
    if r.Method == http.MethodOptions {
        return
    }
	e := &Event{}
	vars := mux.Vars(r)
	user_id, ok := vars["user_id"]

	if !ok {
		http.Error(w, "Missing User ID param", http.StatusBadRequest)
		return
	}
	event_id, ok := vars["event_id"]
	if !ok {
		http.Error(w, "Missing Event ID param", http.StatusBadRequest)
		return
	}
	
	event_id_int, err := strconv.ParseUint(event_id, 10, 32)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Event ID is invalid", http.StatusBadRequest)
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024);
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&e)
	if(err != nil) {
		log.Printf("JSON Decoding error: %v", err)
		http.Error(w, "JSON Decoding Error", http.StatusBadRequest)
	}
	
	e.ID = uint(event_id_int)

	var count int64
	err = s.db.Session(&gorm.Session{}).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&User{}).Where("id = ?", user_id).Count(&count)
		if res.Error != nil {
			return res.Error
		}
		if count == 0 {
			return DoesNotExistError
		}
		
		res = tx.Model(&Event{}).Where("id = ?", e.ID).Updates(&e);

		return nil
	})

	retok := okS{Ok : true}
	retokb, err := json.Marshal(retok)
	if err != nil {
		log.Printf("JSON Marshaling error: %v", err)
		http.Error(w, "JSON Marshaling Error", http.StatusInternalServerError)
	}

	w.Write(retokb)
}

func (s *HttpServer) removeEventByIDUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
    if r.Method == http.MethodOptions {
        return
    }

	e := &Event{}
	vars := mux.Vars(r)
	user_id, ok := vars["user_id"]

	if !ok {
		http.Error(w, "Missing User ID param", http.StatusBadRequest)
		return
	}
	event_id, ok := vars["event_id"]
	if !ok {
		http.Error(w, "Missing Event ID param", http.StatusBadRequest)
		return
	}
	
	var count int64
	err := s.db.Session(&gorm.Session{}).Transaction(func (tx *gorm.DB) error {
		res := tx.Model(&User{}).Where("id = ?", user_id).Count(&count)
		if res.Error != nil {
			return res.Error
		}
		if count == 0 {
			return DoesNotExistError
		}
		res = tx.Model(&Event{}).Where("id = ?", event_id).First(&e)
		if res.Error != nil {
			return res.Error
		}
		res = tx.Model(&Event{}).Delete(&e)

		return nil
	})
	if err != nil {
		log.Printf("GORM Transaction error: %v", err);
		http.Error(w, "GORM TRANSACTION ERROR", http.StatusBadRequest)
	}

	retok := okS{Ok : true}
	retokb, err := json.Marshal(retok)
	if err != nil {
		log.Printf("JSON Marshaling error: %v", err)
		http.Error(w, "JSON Marshaling Error", http.StatusInternalServerError)
	}

	w.Write(retokb)

}

func HttpInitServer(db *gorm.DB, port string) (*HttpServer) {
	h := &HttpServer{};
	h.db = db;
	h.mux = mux.NewRouter()
	h.mux.HandleFunc("/ok", httpOkHandler)
	h.mux.HandleFunc("/auth", h.authHandler).Methods("POST", "OPTIONS")
	h.mux.HandleFunc("/events/add/{user_id}", h.addEventToUser).Methods("POST", "OPTIONS")
	h.mux.HandleFunc("/events/get/{user_id}/all", h.getAllEventsUser).Methods("POST", "OPTIONS")
	h.mux.HandleFunc("/events/get/{user_id}/{event_id}", h.getEventByIDUser).Methods("POST", "OPTIONS")
	h.mux.HandleFunc("/events/get/{user_id}/from/{from}/until/{until}", h.getEventByDateReqAndIDUser).Methods("POST", "OPTIONS")
	h.mux.HandleFunc("/events/set/{user_id}/{event_id}", h.setEventByIDUser).Methods("POST", "OPTIONS")
	h.mux.HandleFunc("/events/rm/{user_id}/{event_id}", h.removeEventByIDUser).Methods("POST", "OPTIONS")
	h.mux.Use(mux.CORSMethodMiddleware(h.mux))

	h.s = &http.Server{
		Addr:	fmt.Sprintf(":%s", port),
		Handler: h.mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	log.Printf("Listening on port :%s", port);

	return h;
}

func (self *HttpServer) HttpListen() error {
	return self.s.ListenAndServe()
}
