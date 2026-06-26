package main

import (
	//"fmt"
	//"os"
	"fmt"
	"os/exec"
	"os"
	//"path/filepath"
	//"bytes"
	"bufio"
	"encoding/json"
	"io"
	"time"
	"strconv"
)

/*
type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
}



func Test() {

	req := Request{
		Name: "bro",
	}

	data, _ := json.Marshal(req)

	cmd := exec.Command(
		"./test.ps1",
	)

	cmd.Stdin = bytes.NewReader(data)

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var res Response

	json.Unmarshal(out, &res)

	fmt.Println(res.Message)
}
*/
/*
	// get absolute paths so pwsh can find everything
	sessionScript, _ := filepath.Abs("session.ps1")
	userScript, _    := filepath.Abs("script.ps1")
    // "--memory", "256m",
	cmd := exec.Command("powershell",
		"-NonInteractive",
		"-NoProfile",
		"-File", sessionScript,
		"-ScriptPath", userScript,
		"-Seed","hhf4",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("running script...")
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("done")
*/



func main(){

	start3 := time.Now()
	ps, err:= NewPSServer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ps.Close()

	fmt.Println("server strated in:", time.Since(start3))


	start5 := time.Now()
	fmt.Println(ps.CreateShortcut("my app",".\\main.exe","--ass --hole"))
	fmt.Println("command run in:", time.Since(start5))


	start7 := time.Now()
	fmt.Println(ps.CreateShortcut("my app 2",".\\sandbox.exe","--ass --hole"))
	fmt.Println("command run in:", time.Since(start7))
}








func Bench() {
	//ps, _ := NewPSServer()
	numberReqest := 30

	fmt.Printf("doing %s rqs\n",strconv.Itoa(numberReqest))



	start3 := time.Now()
	ps, _ := NewPSServer()
	fmt.Println("server strated in:", time.Since(start3))
	

	
	start4 := time.Now()
	s, err := ps.GetDate(2022, 12)
		if err != nil {
			panic(err)
		}
		_ = s



	fmt.Println("first rquest:", time.Since(start4))



	//fmt.Println("slepping for 30s")
	//time.Sleep(30 * time.Second)

	start := time.Now()
	for i := 0; i < numberReqest-1; i++ {
		//fmt.Println("slepping for 1s")
		//time.Sleep(1* time.Second)

		//start := time.Now()
		s, err := ps.GetDate(2022, 12)
		if err != nil {
			panic(err)
		}
		_ = s
		//fmt.Println("req in:", time.Since(start))
	}

	fmt.Println("elapsed (server):", time.Since(start))













	//start2 := time.Now()
	//callExec(numberReqest)
	//fmt.Println("elapsed (exec):", time.Since(start2))
	/*
	start4 := time.Now()
	callExecSingleProcess(numberReqest)
	fmt.Println("elapsed (exec-once/n:cmd):", time.Since(start4))
	*/

	defer ps.Close()


}




func callExec(n int) {
    for i := 0; i < n; i++ {
        cmd := exec.Command(
            "powershell.exe",
            "-NoProfile",
            "-NonInteractive",
            "-Command",  
    		"Get-Date -Year 2022",

        )

        out, _ := cmd.Output()
        _ = out
    }
}


func callExecSingleProcess(n int) {
    script := fmt.Sprintf(`
        for ($i=0; $i -lt %d; $i++) {
            Get-Date -Year 2022
        }
    `, n)

    cmd := exec.Command(
        "powershell.exe",
        "-NoProfile",
        "-NonInteractive",
        "-Command",
        script,
    )

    _, _ = cmd.Output()
}









type Request struct {
    Fid    int `json:"fid"`
	Params map[string]any `json:"params,omitempty"`
}

type Response struct {
    Ok    bool   `json:"ok"`
    Result string `json:"result"`
	Error  string  `json:"error,omitempty"`
}	



type PSServer struct {
    stdin  io.WriteCloser
    reader *bufio.Reader
    cmd    *exec.Cmd
}


func NewPSServer() (*PSServer, error) {
    cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-File", "./server.ps1")

	stdin, err := cmd.StdinPipe()
	if err != nil { return nil, err }

	stdout, err := cmd.StdoutPipe()
	if err != nil { return nil, err }

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	s := &PSServer{
        stdin:  stdin,
        reader: bufio.NewReader(stdout),
        cmd:    cmd,
    }

	
	if err := s.InitServer(); err != nil {
		return nil, err
	}
	


    return s, nil
}



func (s *PSServer) Send(req Request) (Response, error) {
    data, _ := json.Marshal(req)
    io.WriteString(s.stdin, string(data)+"\n")  

    line, _ := s.reader.ReadString('\n')         
	//fmt.Println(line)
    var res Response
    json.Unmarshal([]byte(line), &res)

	if !res.Ok{
		return Response{},fmt.Errorf(res.Error)
	}

    return res, nil
}

func (s *PSServer) NewRequest(fid int, params map[string]any) Request {
    return Request{
        Fid:    fid,
        Params: params,
    }
}

func (s *PSServer) GetDate(year int, month int) (string,error) {
	var fid int = 1 

	req := s.NewRequest(fid,map[string]any{

		"Year": year,
    	"Month": month,

	})


	res ,err := s.Send(req)

	if err != nil{
		return "",err
	}

	return res.Result,nil
}






func (s *PSServer) CreateShortcut(name string, target string, arguments string) (string,error) {
	var fid int = 2 

	params := map[string]any{
		"Name":   name,
		"Target": target,
	}

	if arguments != "" {
    params["Arguments"] = arguments
	}

	req := s.NewRequest(fid, params)

	res ,err := s.Send(req)

	if err != nil{
		return "",err
	}

	return res.Result,nil
}










func (s *PSServer) InitServer() (error) {
	var fid int = 0

	req := s.NewRequest(fid,map[string]any{})


	_ ,err := s.Send(req)

	if err != nil{
		return err
	}


	return nil
}








func (s *PSServer) Close() {
    s.stdin.Close()
    s.cmd.Wait()
}