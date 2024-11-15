package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "strconv"

    "github.com/manifoldco/promptui"
)

type Nomi struct {
    Uuid string
    Gender string
    Name string
    Created string
    RelationshipType string
}

type NomiContainer struct {
    Nomis []Nomi
}

type RoomReceive struct {
    Uuid string
    Name string
    Created string
    Updated string
    Status string
    BackchannelingEnabled bool
    Nomis []Nomi
    Note string
}

type RoomReceiveContainer struct {
    Rooms []RoomReceive
}

type RoomSend struct {
    name string
    note string
    backchannelingEnabled bool
    nomiUuids []string
}

func GetNomiById(nomiCollection []Nomi, uuid string) *Nomi {
    for _, n := range nomiCollection {
        if n.Uuid == uuid {
            return &n
        }
    }

    return nil
}

func GetRoomById(roomCollection []RoomReceive, uuid string) *RoomReceive {
    for _, r := range roomCollection {
        if r.Uuid == uuid {
            return &r
        }
    }

    return nil
}

func (nomi *Nomi) DisplayNomi(mode string) string {
    if mode == "" {
        return fmt.Sprintf(red + "%-20s" + cyan + "%10s\n" + reset, nomi.Name, nomi.Uuid)
    } else {
        if strings.ToLower(mode) == "verbose" {
            return fmt.Sprintf(red + "%v\n" + green + " Uuid: " + cyan + "%v\n" + green + " Gender: %v\n RelationshipType: %v\n Created: %v\n" + reset, nomi.Name, nomi.Uuid, nomi.Gender, nomi.RelationshipType, nomi.Created)
        }
    }
    return ""
}

func (room *RoomReceive) DisplayRoom(mode string) string {
    if mode == "" {
        return fmt.Sprintf(red + "%-20s" + cyan + "%10s\n" + reset, room.Name, room.Uuid)
    } else {
        if strings.ToLower(mode) == "verbose" {
            var retString string
            retString += fmt.Sprintf(yellow + "%v\n" + green + " Uuid: " + cyan + "%v\n" + green + " Created: %v\n Updated: %v\n Status: %v\n Backchanneling: %v\n Note: %v\n Nomis:\n", room.Name, room.Uuid, room.Created, room.Updated, room.Status, room.BackchannelingEnabled, room.Note)
            for _, n := range room.Nomis {
                retString += fmt.Sprint("    ")
                retString += n.DisplayNomi("")
            }
            return retString
        }
    }

    return ""
}

// Color constants
const (
    red     = "\033[31m"
    green   = "\033[32m"
    yellow  = "\033[33m"
    blue    = "\033[34m"
    magenta = "\033[35m"
    cyan    = "\033[36m"
    reset   = "\033[0m"
)

var ApiRoot string
var ApiKey string
var UserNomis []Nomi
var UserRooms []RoomReceive
var promptTemplateNoProps *promptui.SelectTemplates
var promptTemplateWithProps *promptui.SelectTemplates

func main() {
    fmt.Println("Welcome to d3tour's Nomi Room Manager")
    fmt.Println("Use the arrow keys to navigate selection menus")

    ApiRoot = "https://api.nomi.ai/v1/"
    ApiKey = os.Getenv("NOMI_API_KEY")
    var err error
    if ApiKey == "" {
        keyPrompt := promptui.Prompt {
            Label: "Enter your Nomi.ai API key (skip this step in the future by setting your NOMI_API_KEY environment variable)",
        }
        ApiKey, err = keyPrompt.Run()
        if err != nil {
            fmt.Printf("Error entering API key: %v\n", err)
        }
    }

    if ApiKey == "" {
        fmt.Printf("No Nomi.ai API key was provided")
        return
    }

    promptTemplateNoProps = &promptui.SelectTemplates{
        Active: `▶ {{ . | cyan | bold }}`,
        Inactive: `  {{ . | yellow }}`,
        Selected: `✔ {{ . | green | bold }}`,
        Details: `{{ "Selected:" | faint }} {{ . }} `,
    }

    promptTemplateWithProps = &promptui.SelectTemplates{
        Label:    "{{ . }}",
        Active:   "▶ {{ .Name | red }} ({{ .Uuid | cyan }})",
        Inactive: "  {{ .Name | red }} ({{ .Uuid | yellow }})",
        Selected: "✔ {{ .Name | green}} ({{.Uuid | cyan }} selected)",
    }
    stop := false
    for {
        stop = mainMenu()
        if stop {
            return
        }
    }
} 

func mainMenu() bool {
    menuItems := []string{
        "0: Exit",
        "1: List Nomis",
        "2: List Rooms",
        "3: Create Nomi Room",
        "4: Delete Nomi Room",
        "5: Add Nomi To Room",
        "6: Remove Nomi From Room",
        "7: Update Room Name",
        "8: Update Room Note",
        "9: Change Room Backchanneling",
    }

    prompt := promptui.Select{
        Label: "What would you like to do?",
        Items: menuItems,
        Templates: promptTemplateNoProps,
        Size: len(menuItems),
    }

    _, result, err := prompt.Run()
    if err != nil {
        fmt.Printf("Error choosing option: %v\n", err)
    }

    idx := strings.Index(result, ":")
    resultNumber := result[:idx]

    switch resultNumber {
    case "1":
        listNomis(true)
    case "2":
        listRooms(true)
    case "3":
        createRoom()
    case "4":
        deleteRoom()
    case "5":
        addNomiRoom()
    case "6":
        removeNomiRoom()
    case "7":
        updateRoom("name")
    case "8":
        updateRoom("note")
    case "9":
        updateRoom("backchanneling")
    case "0":
        fmt.Println("Bye!")
        return true
    }

    return false
}

func listNomis(display bool) {
    callUrl := strings.Join([]string{ApiRoot, "nomis"}, "")
    apiOut, err := ApiCall(callUrl, "GET", nil)
    if err != nil {
        fmt.Printf("Error in list Nomis API call: %v\n", err)
    }

    var nomis NomiContainer
    if err := json.Unmarshal([]byte(apiOut), &nomis); err != nil {
        fmt.Printf("Error unmarshalling list Nomis API response: %v\n", err)
        return
    }

    UserNomis = nomis.Nomis

    if display {
        for _, n := range UserNomis {
            fmt.Println(n.DisplayNomi("verbose"))
        }
    }
    return
}

func listRooms(display bool) {
    listNomis(false)

    callUrl := strings.Join([]string{ApiRoot, "rooms"}, "")
    apiOut, err := ApiCall(callUrl, "GET", nil)
    if err != nil {
        fmt.Printf("Error in list Nomis API call: %v\n", err)
    }

    var rooms RoomReceiveContainer
    if err := json.Unmarshal([]byte(apiOut), &rooms); err != nil {
        fmt.Printf("Error unmarshalling list rooms API response: %v\n", err)
        return
    }

    UserRooms = rooms.Rooms

    if display {
        for _, r := range UserRooms {
            fmt.Print(r.DisplayRoom("verbose"))
        }
    }
    return
}

func createRoom() {
    listNomis(false)

    var err error
    roomPrompt := promptui.Prompt {
       Label: "Room Name",
    }

    roomName, err := roomPrompt.Run()
    if err != nil {
        fmt.Printf("Error getting room name: %v\n", err)
    }

    roomPrompt.Label = "Room Note"
    roomNote, err := roomPrompt.Run()
    if err != nil {
        fmt.Printf("Error getting room name: %v\n", err)
    }

    backchannelPrompt := promptui.Select{
        Label: "Room Backchanneling",
        Items: []bool{true, false},
        Templates: promptTemplateNoProps,
    }

    _, backchanneling, err := backchannelPrompt.Run()
    if err != nil {
        fmt.Printf("Error choosing backchanneling option: %v\n", err)
    }
    backchannelBool, err := strconv.ParseBool(backchanneling)
    if err != nil {
        fmt.Printf("Error parsing backchanneling option as bool: %v", err)
    }

    nomisToAdd := nomiMultiSelect(UserNomis)
    var nomiUuids []string
    for _, n := range nomisToAdd {
        nomiUuids = append(nomiUuids, n.Uuid)
    }

    roomToCreate := map[string]interface{}{
        "name": roomName,
        "note": roomNote,
        "backchannelingEnabled": backchannelBool,
        "nomiUuids": nomiUuids,
    }

    callUrl := strings.Join([]string{ApiRoot, "rooms"}, "")
    apiOut, err := ApiCall(callUrl, "POST", roomToCreate)
    if err != nil {
        fmt.Printf("Error in list Nomis API call: %v\n", err)
    }

    var room RoomReceive
    if err := json.Unmarshal([]byte(apiOut), &room); err != nil {
        fmt.Printf("Error unmarshalling list rooms API response: %v\n", err)
        return
    }

    fmt.Println("Created Room:")
    fmt.Println(room.DisplayRoom("verbose"))
}

func nomiMultiSelect(nomis []Nomi) []Nomi {
    var retItems []Nomi
    var choices []interface{}
    choices = append(choices, Nomi{
        Name: "Finish Selection",
        Uuid: "finished",
        Gender: " ",
        Created: " ",
        RelationshipType: " ",
    })
    for _, n := range nomis {
        choices = append(choices, n)
    }

    selectedItems := make(map[string]Nomi)

    for {
        prompt := promptui.Select{
            Templates: promptTemplateWithProps,
            Label: "Select Nomis (Press enter to toggle selection. Select 'Finish Selection' to end.)",
            Items: choices,
            Size: len(nomis) + 1,
        }

        _, choice, err := prompt.Run()
        if err != nil {
            fmt.Printf("Error choosing Nomis: %v\n", err)
        }

        if strings.TrimPrefix(strings.Split(choice, " ")[0], "{") == "finished" {
            break
        }

        selectedItem := GetNomiById(UserNomis, strings.TrimPrefix(strings.Split(choice, " ")[0], "{"))

        if _, exists := selectedItems[selectedItem.Uuid]; exists {
            delete(selectedItems, selectedItem.Uuid)
        } else {
            selectedItems[selectedItem.Uuid] = *selectedItem
        }
    }

    for _, n := range selectedItems {
        retItems = append(retItems, n)
    }

    return retItems
}

func deleteRoom() {
    listRooms(false)

    deleteRoomPrompt := promptui.Select{
        Label: "Choose a room to delete",
        Items: UserRooms,
        Templates: promptTemplateWithProps,
        Size: len(UserRooms),
    }

    _, roomToDelete, err := deleteRoomPrompt.Run()
    if err != nil {
        fmt.Printf("Error choosing room to delete option: %v\n", err)
    }

    callUrl := strings.Join([]string{ApiRoot, "rooms/", strings.TrimPrefix(strings.Split(roomToDelete, " ")[0], "{")}, "")
    apiOut, err := ApiCall(callUrl, "DELETE", nil)
    if err != nil {
        fmt.Printf("Error in list Nomis API call: %v\n", err)
    }

    roomDeleteResponse := string(apiOut)

    if strings.TrimSpace(roomDeleteResponse) == "" {
        fmt.Printf("Deleted room: %v (%v)\n", strings.TrimPrefix(strings.Split(roomToDelete, " ")[0], "{"), strings.TrimPrefix(strings.Split(roomToDelete, " ")[2], "{"))
    } else {
        fmt.Printf("Delete room sent back a response which is a bad thing:\n %v\n", roomDeleteResponse)
    }

    listRooms(false)

    fmt.Println("\nCurrent Rooms:")
    for _, r := range UserRooms {
        fmt.Print(r.DisplayRoom(""))
    }

    return
}

func addNomiRoom() {
    listRooms(false)

    addNomiRoomPrompt := promptui.Select{
        Label: "Choose a room to add a Nomi to",
        Items: UserRooms,
        Templates: promptTemplateWithProps,
        Size: len(UserRooms),
    }

    _, roomToAddTo, err := addNomiRoomPrompt.Run()
    if err != nil {
        fmt.Printf("Error choosing room to add Nomi to: %v\n", err)
    }

    roomUuid := strings.TrimPrefix(strings.Split(roomToAddTo, " ")[0], "{") 
    callUrl := strings.Join([]string{ApiRoot, "rooms/", roomUuid}, "")
    nomisToAdd := nomiMultiSelect(UserNomis)
    var nomisToAddUuids []string
    for _, n := range nomisToAdd {
        nomisToAddUuids = append(nomisToAddUuids, n.Uuid)
    }
    for _, n := range GetRoomById(UserRooms, roomUuid).Nomis {
        nomisToAddUuids = append(nomisToAddUuids, n.Uuid)
    }

    callBody := map[string]interface{}{
        "nomiUuids": nomisToAddUuids,
    }

    _, err = ApiCall(callUrl, "PUT", callBody)
    if err != nil {
        fmt.Printf("Error in add Nomi to room API call: %v\n", err)
    }

    listRooms(false)
    fmt.Println("Updated room:")
    roomUpdated := GetRoomById(UserRooms, roomUuid)
    fmt.Println(roomUpdated.DisplayRoom("verbose"))

    return
}

func removeNomiRoom() {
    listRooms(false)

    removeNomiRoomPrompt := promptui.Select{
        Label: "Choose a room to remove a Nomi from",
        Items: UserRooms,
        Templates: promptTemplateWithProps,
        Size: len(UserRooms),
    }

    _, roomToRemoveFrom, err := removeNomiRoomPrompt.Run()
    if err != nil {
        fmt.Printf("Error choosing room to add Nomi to: %v\n", err)
    }

    roomUuid := strings.TrimPrefix(strings.Split(roomToRemoveFrom, " ")[0], "{") 
    callUrl := strings.Join([]string{ApiRoot, "rooms/", roomUuid}, "")
    roomNomis := GetRoomById(UserRooms, roomUuid).Nomis
    nomisToRemove := nomiMultiSelect(roomNomis)

    var nomisToRemoveUuids []string
    for _, n := range nomisToRemove {
        nomisToRemoveUuids = append(nomisToRemoveUuids, n.Uuid)
    }

    var currentRoomNomiUuids []string
    for _, n := range GetRoomById(UserRooms, roomUuid).Nomis {
        currentRoomNomiUuids = append(currentRoomNomiUuids, n.Uuid)
    }

    resultantNomiUuids := removeFromSlice(currentRoomNomiUuids, nomisToRemoveUuids)
    fmt.Printf("Resultant IDs: %v\n", strings.Join(resultantNomiUuids, ", "))

    callBody := map[string]interface{}{
        "nomiUuids": resultantNomiUuids,
    }

    _, err = ApiCall(callUrl, "PUT", callBody)
    if err != nil {
        fmt.Printf("Error in add Nomi to room API call: %v\n", err)
    }

    listRooms(false)
    fmt.Println("Updated room:")
    roomUpdated := GetRoomById(UserRooms, roomUuid)
    fmt.Println(roomUpdated.DisplayRoom("verbose"))

    return
}

func removeFromSlice(original []string, toRemove []string) []string {
    removeSet := make(map[string]bool)
    for _, i := range toRemove {
        removeSet[i] = true
    }

    retSlice := []string{}
    for _, j := range original {
        if !removeSet[j] {
            retSlice = append(retSlice, j)
        }
    }

    return retSlice
}

func updateRoom(property string) {
    listRooms(false)

    roomUpdatePrompt := promptui.Select{
        Label: "Choose a room to update",
        Items: UserRooms,
        Templates: promptTemplateWithProps,
        Size: len(UserRooms),
    }

    _, roomToUpdateString, err := roomUpdatePrompt.Run()
    if err != nil {
        fmt.Printf("Error choosing room to add Nomi to: %v\n", err)
    }


    roomUuid := strings.TrimPrefix(strings.Split(roomToUpdateString, " ")[0], "{") 
    roomToUpdate := GetRoomById(UserRooms, roomUuid)
    
    labelText := ""
    currentValue := ""
    jsonBodyProperty := ""
    var updateValue interface{}

    switch strings.ToUpper(property) {
        case "NAME":
        labelText = "Name"
        currentValue = roomToUpdate.Name
        jsonBodyProperty = "name"
        case "NOTE":
        labelText = "Note"
        currentValue = roomToUpdate.Note
        jsonBodyProperty = "note"
        case "BACKCHANNELING":
        labelText = "Backchanneling Setting"
        currentValue = strconv.FormatBool(roomToUpdate.BackchannelingEnabled)
        jsonBodyProperty = "backchannelingEnabled"
    }

    promptLabel := "Update " + labelText + " (Current Value: " + currentValue + ")"

    if jsonBodyProperty == "name" || jsonBodyProperty == "note" {
        roomPrompt := promptui.Prompt {
            Label: promptLabel,
        }

        updateValue, err = roomPrompt.Run()
        if err != nil {
            fmt.Printf("Error getting new room value: %v\n", err)
        }
    } else {
        backchannelPrompt := promptui.Select{
            Label: promptLabel,
            Items: []bool{true, false},
            Templates: promptTemplateNoProps,
        }

        _, bcValue, err := backchannelPrompt.Run()
        if err != nil {
            fmt.Printf("Error choosing backchanneling option: %v\n", err)
        }

        updateValue, err = strconv.ParseBool(bcValue)
        if err != nil {
            fmt.Printf("Error converting selection %v to bool: %v", bcValue, err)
        }
    }

    callBody := map[string]interface{}{
        jsonBodyProperty: updateValue,
    }
    callUrl := strings.Join([]string{ApiRoot, "rooms/", roomUuid}, "")

    _, err = ApiCall(callUrl, "PUT", callBody)
    if err != nil {
        fmt.Printf("Error in add Nomi to room API call: %v\n", err)
    }

    listRooms(false)
    fmt.Println("Updated room:")
    roomUpdated := GetRoomById(UserRooms, roomUuid)
    fmt.Println(roomUpdated.DisplayRoom("verbose"))

    return
}

func ApiCall(endpoint string, method string, body interface{}) ([]byte, error) {
    method = strings.ToUpper(method)

    var jsonBody []byte
    var bodyReader io.Reader
    var err error

    if body != nil {
        jsonBody, err = json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("Error constructing body: %v: ", err)
        }
        bodyReader = bytes.NewBuffer(jsonBody)
    } else {
        bodyReader = nil
    }

    req, err := http.NewRequest(method, endpoint, bodyReader)
    if err != nil {
        return nil, fmt.Errorf("Error reading HTTP request: %v", err)
    }

    req.Header.Set("Authorization", ApiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("Error making HTTP request: %v", err)
    }

    defer resp.Body.Close()

    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("Error reading HTTP response: %v", err)
    }

    if resp.StatusCode < 200 || resp.StatusCode > 299 {
        var errorResult map[string]interface{}
        if err := json.Unmarshal(responseBody, &errorResult); err != nil {
            return nil, fmt.Errorf("Error unmarshalling API error response: %v\n%v", err, string(responseBody))
        }
        return nil, fmt.Errorf("Error response from Nomi API\n Error Code: %v\n Response Body: %v\n",resp.StatusCode, string(responseBody))
    }

    return responseBody, nil
}
