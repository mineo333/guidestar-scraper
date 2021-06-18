package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "io"
  "os"
  "golang.org/x/net/html"
  "strings"
  "encoding/json"
  "net/url"
  "time"
  "math/rand"
)
type program struct{
    Name string
    Description string
}
type nonprofit struct{
  Name string
  Ein string
  Mission_statement string
  Programs []program
}
var base_url string = "https://www.guidestar.org/profile/" //the base url that will be used


func get_html_for_ein(ein string) io.ReadCloser{ //this gets the html for an ein as an io reader because the html parser requires an ioreader
  url := base_url + ein;
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return nil
  }
  req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36") //weird shit. the website only requires a browser user-agent and no authentication cookies. ill take it tho. This header was obtained by copying a google chrome header from dev mode

  res, err := client.Do(req)
  //fmt.Println(res.StatusCode)
  if(res.StatusCode != 200){
    return nil // could not be found. This is usually 404.
  }
  if err != nil {
    fmt.Println(err)
    return nil
  }


  return res.Body;
}


func getPrograms(z *html.Tokenizer) []program{ //returns list of program structs and size
  programs := make([]program, 0) // i don't know how many programs there are so we're doing an arraylist type thing here
  var val []byte
  _, val, _ =  z.TagAttr() //get current val
  for !strings.Contains(string(val), "theMaps") {  //theMaps directly comes after the programs

    _, val, _ =  z.TagAttr()
    if strings.Contains(string(val), "programHead") {
      var pinfo program //final struct
       z.Next() //this is hardcoded and is standardized on the site.
       z.Next()
       z.Next()
       z.Next()
       z.Next()
       z.Next()
      pinfo.Name = strings.TrimSpace(string(z.Text()))
      for !strings.Contains(string(z.Raw()),"description"){
        z.Next()
      }
      z.Next()
      pinfo.Description = strings.TrimSpace(string(z.Text()))
      programs = append(programs, pinfo)
    }

    z.Next()
  }





  return programs
}
func getNonProfit(ein string, name string) *nonprofit{
  var mission_statement string;
  var programs []program
  var nonprofit nonprofit
  body := get_html_for_ein(ein);
  if(body == nil){ //the request didn't go through properly
    fmt.Printf("ein: %s was not found\n", ein)
    return nil;
  }
  z := html.NewTokenizer(body)

  for /*i:=0; i<1000; i++*/z.Err() == nil{
      _, val, _ :=  z.TagAttr()

    if strings.Contains(string(val), "mission-statement"){ //gets mission statement. Mission statement always precedes programs
      z.Next()
      mission_statement = string(z.Text()) //NOTE: for error checking keep in mind that this can sometimes be "This organization has not provided GuideStar with a mission statement.". Examples of broken ones are 01-0581159 and 01-0574855
    }
    if strings.Contains(string(val), "progamsAccordion"){ //we're now ready to begin parsing the program list. If there is no Accordion there is no list of programs
      //fmt.Println("Got Here")

      programs = getPrograms(z)
      break
    }
    z.Next()

  }
  nonprofit.Name = name
  nonprofit.Mission_statement = mission_statement
  nonprofit.Ein = ein
  nonprofit.Programs = programs
  body.Close();
  return &nonprofit
}

func getJson(np nonprofit) []byte{
  bytes, err:= json.Marshal(np)
  //fmt.Println(np)
  if err != nil{
    fmt.Println("error")
    return nil
  }
  return bytes
}

func getEin(np_name string) string{ // for getting the ein i will only use the first result in the nonprofit array because the rest are typically redundant and only represent alternate locations
  enc_np_name := url.QueryEscape(np_name)

  var init interface{}
  url := "https://projects.propublica.org/nonprofits/api/v2/search.json?q="+enc_np_name
  method := "GET"

  client := &http.Client {}
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return ""
  }
  res, err := client.Do(req)
  if res.StatusCode != 200{
    fmt.Print(res.StatusCode)
    fmt.Print(" ")
    if res.StatusCode == 403{
      return "403"
    }
    return ""
  }
  if err != nil {
    fmt.Println(err)
    return ""
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return ""
  }
  json.Unmarshal(body, &init)
  s1 ,ok := init.(map[string]interface{})["organizations"].([]interface{})
  if(!ok){
    return ""
  }

  ein := s1[0].(map[string]interface{})["strein"].(string)
  return ein





}

func main() {


  var eins []string = make([]string, 0)
  nonprofits := make([]nonprofit, 0)



  name_file, err := os.Open("npnames")
  if err != nil{
    fmt.Println("File does not exist")
    return
  }
  stat_file,_ := name_file.Stat()
  file_len := stat_file.Size()
  buf := make([]byte, file_len)//will contain file_
  _, err = name_file.Read(buf)
  if err != nil{
    fmt.Println("Something went wrong when reading the file")
  }
  buf_string := string(buf)
  np_strings := strings.Split(buf_string, "\n")
  //parsing of file is complete



  fd, _ := os.OpenFile("missing_eins", os.O_CREATE | os.O_RDWR, 0755)
  fd.Seek(0,2)
  fmt.Println("Note if the error code at the beginning does not say 404 then the site is blocking the search.")
  for _, str := range np_strings{
    tts := time.Duration(rand.Intn(500-400)+400)
    time.Sleep(tts*time.Millisecond) //random sleep time do i don't get BLOCKED
    ein := getEin(strings.TrimSpace(str))
    for ein == "403"{
      fmt.Println("got a code 403. Sleeping for 60 seconds to reset")
      time.Sleep(60*time.Second) //reset real quick
      ein = getEin(strings.TrimSpace(str)) //get it again
    }
    if ein == ""{
      fmt.Println("EIN could not be found for: " + str + " This name will be added to the \"missing_ein's\" file")
      fd.Write([]byte(str + "\n"))
      //we shouldn't continue here so that in the next for loop the counter, i, is accurate. We want the eins array and the np_strings array to be the same size
    }
    eins = append(eins, ein)
  }







  fmt.Println("Getting nonprofit info from GuideStar")

  for i, ein:= range eins{ //iterate through each ein provided
    tts := time.Duration(rand.Intn(500-400)+400)
    time.Sleep(tts*time.Millisecond) //random sleep time do i don't get BLOCKED
    if ein == ""{
      continue //this is an error ein
    }
    np := getNonProfit(ein, np_strings[i])
    if np == nil{ //we got an error :(404 EIN could not be found for: actionaid international italia This name will be added to the "missing_ein's" file

      fmt.Println("The information for "+np_strings[i]+" could not be found by GuideStar or there was an error\n")
      fd.Write([]byte(np_strings[i] + "\ns")) //write to missing_eins
      continue
    }
    nonprofits = append(nonprofits,*np) //REMINDER: Do some error checking here
  }
  fmt.Println("Writing to files")
  fd.Close()
  for _, np := range nonprofits{
    bytes := getJson(np) //its now time to write the json to a file
    if bytes == nil{
      bytes = []byte("Could something went wrong for " + np.Ein)
    }

    fd, err = os.OpenFile("nonprofit_info/" + np.Name, os.O_CREATE | os.O_RDWR, 0755)

    fd.Write(bytes)
    fd.Close()
  }

}
