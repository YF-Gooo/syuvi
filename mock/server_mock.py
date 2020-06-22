import requests
if __name__ == '__main__':
    data={"target":"123","content":"helloworld"}
    requests.post("http://localhost:8080/command",data)