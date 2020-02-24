package controller

//
//func TestArticleController_GetAll(t *testing.T) {
//	articleHandler := newArticle(nil)
//	req, err := http.NewRequest("GET", "localhost:8000/", nil)
//	//GetAll
//	if err != nil {
//		t.Fatalf("could not create request %v", err)
//	}
//	rec := httptest.NewRecorder()
//	articleHandler.GetAll(rec, req)
//
//	res := rec.Result()
//	if res.StatusCode != http.StatusOK {
//		t.Errorf("expected status OK; got %v", res.Status)
//	}
//
//	b, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		t.Fatalf("could not read response: %v", err)
//	}
//	_, err = strconv.Atoi(string(b))
//	if err != nil {
//		t.Fatalf("expected an integer; got %s,",err.Error())
//	}
//
//
//
//}
