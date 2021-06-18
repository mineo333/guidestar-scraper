package main
import (
  "fmt"
  "encoding/json"
  "github.com/360EntSecGroup-Skylar/excelize/v2"
  "os"
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
func readFile(fname string, size int64) []byte{
  fd, err := os.Open("../nonprofit_info/"+fname)
  defer fd.Close()
  buf := make([]byte, size)


  if err != nil{
    fmt.Println("File " + fname + " does not exist")
    return nil
  }
  _, err = fd.Read(buf)
  if err != nil{
    fmt.Println("Read failed")
    return nil
  }
  return buf
}

func main(){
  fmt.Println("Beginning JSON to Excel")
  spreadsheet := excelize.NewFile()
  main_sheet := spreadsheet.NewSheet("nonprofits")
  spreadsheet.SetActiveSheet(main_sheet)
  //sheet has been created
  files, err := os.ReadDir("../nonprofit_info") //there should be no dir in nonprofit_info
  if err != nil{
    fmt.Println("Error with reading directory")
    return
  }
  for i, file := range files{
    fInfo, _ := file.Info()
    var nonprofit nonprofit
    buf := readFile(file.Name(), fInfo.Size())
    json.Unmarshal(buf, &nonprofit)

    name_loc := fmt.Sprintf("%s%d", "A", i+1) //the reason this is i+1 is because i is 0-indexed. So, to get the proper row we need to
    ein_loc := fmt.Sprintf("%s%d", "B", i+1)
    ms_loc := fmt.Sprintf("%s%d", "C", i+1)
    spreadsheet.SetCellValue("nonprofits", name_loc, nonprofit.Name)
    spreadsheet.SetCellValue("nonprofits", ein_loc, nonprofit.Ein)
    spreadsheet.SetCellValue("nonprofits", ms_loc, nonprofit.Mission_statement)
    //it is time to iterate through the programs of the nonprofit
    if nonprofit.Programs == nil{
      continue //it can be nil at times
    }
    counter := int('C')+1
    for _, prog := range nonprofit.Programs{
      prog_name_loc := fmt.Sprintf("%c%d", rune(counter), i+1)
      desc_name_loc := fmt.Sprintf("%c%d", rune(counter+1), i+1)
      spreadsheet.SetCellValue("nonprofits", prog_name_loc, prog.Name)
      spreadsheet.SetCellValue("nonprofits", desc_name_loc, prog.Description)
      counter += 2
    }

  }
  spreadsheet.SaveAs("nonprofits.xlsx")




}
