#ifndef SOCKETUNIXSERVER_H
#define SOCKETUNIXSERVER_H
#include <sys/un.h>

class socketUnixServer
{
public:
    socketUnixServer(const char * path);
    void send(const char * data);
    void receive();
    virtual void routeHand(const char *data)=0;
    void ReqConfig();

private:
    char * path;
    sockaddr_un addr;
    sockaddr_un addrTo;
};

#endif // SOCKETUNIXSERVER_H
