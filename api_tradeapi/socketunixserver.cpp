#include "socketunixserver.h"
#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <string.h>
#include <iostream>


using namespace std;
socketUnixServer::socketUnixServer(const char * path)
{
    //strcat(p,path);
    //cout<<path<<endl;
    unlink(path);
    this->addrTo.sun_family = AF_UNIX;
    strcpy(this->addrTo.sun_path, path);
    char p[1024];
    strcpy(p,path);
    strcat(p,"_");
    //cout<<p<<endl;
    this->addr.sun_family = AF_UNIX;
    strcpy(this->addr.sun_path, p);
    //this->over = false;
    //this->receive();
}
void socketUnixServer::send(const char *data){
    cout<<data<<endl;
    return;
    int sock;
    //sockaddr_un addr;
    //socklen_t addrlen;
    sock = socket(AF_UNIX, SOCK_DGRAM, 0);
    sendto(sock, data, strlen(data), 0, (sockaddr*)&this->addr, sizeof(this->addr));
    close(sock);

}
void socketUnixServer::receive(){
    int sock;
    //sockaddr_un addr;
    //socklen_t addrlen;
    char buf[1024];
    int n;

    sock = socket(AF_UNIX, SOCK_DGRAM, 0);
    //addr.sun_family = AF_UNIX;
    //strcpy(addr.sun_path, "/tmp/afu_dgram");
    bind(sock, (sockaddr*)&this->addrTo, sizeof(this->addrTo));
    while(true){
      memset(buf, 0, sizeof(buf));
      n = recv(sock, buf, sizeof(buf) - 1, 0);
      this->routeHand(buf);
      //printf("recv:%s\n", buf);
    }
    printf("end socket server");
    close(sock);
}
void socketUnixServer::ReqConfig(){
    this->send("config");
}
