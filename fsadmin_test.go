package gowfs

import "fmt"
import "log"
import "net/url"
import "net/http"
import "net/http/httptest"

import "testing"


func Test_Rename(t *testing.T){
	server := mockServerFor_Rename()
	defer server.Close()
	t.Logf("Started httptest.Server on %v", server.URL)

	url,_ := url.Parse(server.URL)
	conf := Configuration{Addr: url.Host }
	fs, _ := NewFileSystem(conf)
	
	ok, err := fs.Rename(Path{Path:"/testing"}, Path{Path:"/testing/newname"})
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("Rename() - is not renaming value properly")
	}
}

func Test_MkDirs(t *testing.T){
	server := mockServerFor_MkDirs()
	defer server.Close()
	t.Logf("Started httptest.Server on %v", server.URL)

	url,_ := url.Parse(server.URL)
	conf := Configuration{Addr: url.Host }
	fs, _ := NewFileSystem(conf)
	
	ok, err := fs.MkDirs(Path{Path:"/test"}, 0744)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("MkDirs() - is not returning expected FileStatus value")
	}
}

func Test_CreateSymlink(t *testing.T) {
	server := mockServerFor_Symlink()
	defer server.Close()
	t.Logf("Started httptest.Server on %v", server.URL)

	url,_ := url.Parse(server.URL)
	conf  := Configuration{Addr: url.Host }
	fs, _ := NewFileSystem(conf)

	ok, err := fs.CreateSymlink(Path{Path:"/test/orig"}, Path{Path:"/symlink"}, false)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("MkDirs - is not returning expected FileStatus value")
	}
}

func Test_GetFileStatus(t *testing.T){
	server := mockServerFor_FileStatus()
	defer server.Close()
	t.Logf("Started httptest.Server on %v", server.URL)

	url,_ := url.Parse(server.URL)

	conf := Configuration{Addr: url.Host }
	fs, _ := NewFileSystem(conf)
	
	fileStatus, err := fs.GetFileStatus(Path{Path:"/test"})
	if err != nil {
		t.Fatal(err)
	}

	if fileStatus.Permission != "777" || fileStatus.Type == "DIRECORY" {
		t.Fatal("GetFileStatus - is not returning expected FileStatus value")
	}
}

func Test_ListStatus(t *testing.T){
	server := mockServerFor_ListStatus()
	defer server.Close()
	t.Logf("Started httptest.Server on %v", server.URL)

	url,_ := url.Parse(server.URL)

	conf := Configuration{Addr: url.Host }
	fs, _ := NewFileSystem(conf)
	
	statuses, err := fs.ListStatus(Path{Path:"/test"})
	if err != nil {
		t.Fatal(err)
	}

	if len(statuses) != 2 {
		t.Errorf("ListStatus - expecting %d items, but got %d.", 2, len(statuses))
	}
}

func Test_GetContentSummary(t *testing.T) {
	server := mockServerFor_ContentSummary()
	defer server.Close()
	t.Logf("Started httptest.Server on %v", server.URL)

	url, _ := url.Parse(server.URL)
	fs, _  := NewFileSystem(Configuration{Addr:url.Host})
	summary, err := fs.GetContentSummary(Path{Path:"/test"})
	if err != nil {
		t.Fatal (err)
	}
	if summary.SpaceConsumed != 24930 {
		t.Errorf("GetContentSummary - not returning expected values <<%v>>", summary)
	}
}

func Test_GetFileChecksum(t *testing.T) {
	server := mockServerFor_FileChecksum()
	defer server.Close()
	t.Logf("Started httptest.Server on %v", server.URL)

	url, _ := url.Parse(server.URL)
	fs, _  := NewFileSystem(Configuration{Addr:url.Host})
	checksum, err := fs.GetFileChecksum(Path{Path:"/test"})
	if err != nil {
		t.Fatal (err)
	}
	if checksum.Algorithm != "MD5-of-1MD5-of-512CRC32" ||
	   checksum.Length != 28 {
		t.Errorf("GetFileChecksum - not returning expected values <<%v>>", checksum)
	}
}

// *********************** Mock Servers ********************* //
func mockServerFor_Rename() *httptest.Server{
  handler := func (rsp http.ResponseWriter, req *http.Request){
  	  if req.Method != "PUT"{
  	      log.Fatalf("Expecting Request.Method PUT, but got %v", req.Method)
  	  }  	

      q := req.URL.Query()
      if q.Get("op") != OP_RENAME {
        log.Fatalf("Server Missing expected URL parameter: op= %v", OP_RENAME)
      }
      if q.Get("destination") != "/testing/newname" {
          log.Fatalf("Expected param destination to be /testing/newname, but was %v", q.Get("destination"))
      }

      fmt.Fprintf (rsp, `{"Boolean":true}`)
  }
  return httptest.NewServer(http.HandlerFunc(handler))
}

func mockServerFor_MkDirs() *httptest.Server{
  handler := func (rsp http.ResponseWriter, req *http.Request){
  	  if req.Method != "PUT"{
  	      log.Fatalf("Expecting Request.Method PUT, but got %v", req.Method)
  	  }  	
      q := req.URL.Query()
      if q.Get("op") != OP_MKDIRS {
        log.Fatalf("Server Missing expected URL parameter: op= %v", OP_MKDIRS)
      }
      if q.Get("permission") != "744" {
          log.Fatalf("Expected param permission to be 744, but was %v", q.Get("permission"))
      }

      fmt.Fprintf (rsp, `{"Boolean":true}`)
  }
  return httptest.NewServer(http.HandlerFunc(handler))
}


func mockServerFor_Symlink() *httptest.Server{
  handler := func (rsp http.ResponseWriter, req *http.Request){
  	  if req.Method != "PUT"{
  	      log.Fatalf("Expecting Request.Method PUT, but got %v", req.Method)
  	  }  	
      q := req.URL.Query()
      if q.Get("op") != OP_CREATESYMLINK {
        log.Fatalf("Server Missing expected URL parameter: op= %v", OP_CREATESYMLINK)
      }
      if q.Get("destination") != "/test/orig" {
          log.Fatalf("Expected param destination to be /test/orig, but was %v", q.Get("destination"))
      }
      if q.Get("createParent") != "false" {
          log.Fatalf("Expected param createParent to be false, but was %v", q.Get("createParent"))
      }

      fmt.Fprintf (rsp, "")
  }
  return httptest.NewServer(http.HandlerFunc(handler))
}


const listStatusRsp =`
{
  "FileStatuses":
  {
    "FileStatus":
    [
      {
        "accessTime"      : 1320171722771,
        "blockSize"       : 33554432,
        "group"           : "supergroup",
        "length"          : 24930,
        "modificationTime": 1320171722771,
        "owner"           : "webuser",
        "pathSuffix"      : "a.patch",
        "permission"      : "644",
        "replication"     : 1,
        "type"            : "FILE"
      },
      {
        "accessTime"      : 0,
        "blockSize"       : 0,
        "group"           : "supergroup",
        "length"          : 0,
        "modificationTime": 1320895981256,
        "owner"           : "szetszwo",
        "pathSuffix"      : "bar",
        "permission"      : "711",
        "replication"     : 0,
        "type"            : "DIRECTORY"
      }
    ]
  }
}
`
func mockServerFor_ListStatus() *httptest.Server {
	handler := func (rsp http.ResponseWriter, req *http.Request){
  		if req.Method != "GET"{
  	      log.Fatalf("Expecting Request.Method GET, but got %v", req.Method)
  	  	}  	
		q := req.URL.Query()
		if q.Get("op") != OP_LISTSTATUS {
			panic (`Server Missing expected URL parameter: op=` + OP_LISTSTATUS)
		}
		fmt.Fprintln (rsp, listStatusRsp)
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}

const fileStatusRsp =`
{
  "FileStatus":
  {
    "accessTime"      : 0,
    "blockSize"       : 0,
    "group"           : "supergroup",
    "length"          : 0,
    "modificationTime": 1320173277227,
    "owner"           : "webuser",
    "pathSuffix"      : "",
    "permission"      : "777",
    "replication"     : 0,
    "type"            : "DIRECTORY" 
  }
}
`
func mockServerFor_FileStatus() *httptest.Server {
	handler := func (rsp http.ResponseWriter, req *http.Request){
  	  	if req.Method != "GET"{
  	    	log.Fatalf("Expecting Request.Method GET, but got %v", req.Method)
  	  	}  	
		q := req.URL.Query()
		if q.Get("op") != OP_GETFILESTATUS {
			panic (`Server Missing expected URL parameter: op=` + OP_GETFILESTATUS)
		}
		fmt.Fprintln (rsp, fileStatusRsp)
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}

const contentSummaryRsp =`
{
  "ContentSummary":
  {
    "directoryCount": 2,
    "fileCount"     : 1,
    "length"        : 24930,
    "quota"         : -1,
    "spaceConsumed" : 24930,
    "spaceQuota"    : -1
  }
}
`
func mockServerFor_ContentSummary() *httptest.Server {
  handler := func (rsp http.ResponseWriter, req *http.Request){

    q := req.URL.Query()
    if q.Get("op") != OP_GETCONTENTSUMMARY {
      panic (`Server Missing expected URL parameter: op=` + OP_GETCONTENTSUMMARY)
    }
    fmt.Fprintln (rsp, contentSummaryRsp)
  }
  return httptest.NewServer(http.HandlerFunc(handler))
}

const fileChecksumRsp = `
{
  "FileChecksum":
  {
    "algorithm": "MD5-of-1MD5-of-512CRC32",
    "bytes"    : "eadb10de24aa315748930df6e185c0d ...",
    "length"   : 28
  }
}
`
func mockServerFor_FileChecksum() *httptest.Server {
  handler := func (rsp http.ResponseWriter, req *http.Request){
  	if req.Method != "GET"{
  		log.Fatalf("Expecting Request.Method GET, but got %v", req.Method)
  	}

  	q := req.URL.Query()
    if q.Get("op") != OP_GETFILECHECKSUM {
      panic (`Server Missing expected URL parameter: op=` + OP_GETFILECHECKSUM)
    }
    fmt.Fprintln (rsp, fileChecksumRsp)
  }
  return httptest.NewServer(http.HandlerFunc(handler))	
}
