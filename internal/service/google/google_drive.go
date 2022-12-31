package google

var (
//	googleOauthConfig = &oauth2.Config{
//		RedirectURL:  "http://localhost:3000/GoogleCallback",
//		ClientID:     os.Getenv("googlekey"),    // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
//		ClientSecret: os.Getenv("googlesecret"), // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
//		Scopes:       []string{"https://www.googleapis.com/auth/drive", "https://www.googleapis.com/auth/drive.file"},
//		Endpoint:     google.Endpoint,
//	}
//
// Some random string, random for each request
// oauthStateString = "random"
)

//func Conect() {
//
//	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
//	driveService, err := drive.New(client)
//	if err != nil {
//		fmt.Fprintln(w, err)
//		return
//	}
//
//}

//func handleGoogleLogin() {
//	url := googleOauthConfig.AuthCodeURL(oauthStateString)
//	http.Redirect(w, r, url, http.StatusTemporaryRedirect) //delete w and r
//}

//func handleGoogleCallback() {
//	//state := r.FormValue("state")
//	//if state != oauthStateString {
//	//	fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
//	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
//	//	return
//	//}
//
//	//code := r.FormValue("code")
//	token, err := googleOauthConfig.Exchange(context.Background(), code)
//	if err != nil {
//		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
//		return
//	}
//
//	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
//
//}

//func Update(name, path, login, password string) error {
//	Connect(path, login, password)
//	if err != nil {
//		return err
//	}
//	myFile := drive.File{Name: name}
//	driveService, err := drive.New(client)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	createdFile, err := driveService.Files.Create(&myFile).Do()
//	if err != nil {
//		log.Println("Couldn't create file ", err)
//		return err
//	}
//
//	myQR, err := os.Open(path)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	updatedFile := drive.File{Name: name}
//
//	driveService.Files.Update(createdFile.Id, &updatedFile)
//
//	_, err = driveService.Files.Update(createdFile.Id, &updatedFile).Media(myQR).Do()
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	return nil
//
//}
