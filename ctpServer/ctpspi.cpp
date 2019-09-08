#include <iostream>
#include "ctpspi.h"
#include <string.h>
#include <thread>
using namespace std;

ctpspi::ctpspi(const char *BrokerID, const char *UserID, const char *Password)
{
    memset(&this->UserReq,0,sizeof(this->UserReq));
    strcpy(this->UserReq.BrokerID,BrokerID);
    strcpy(this->UserReq.UserID,UserID);
    strcpy(this->UserReq.Password,Password);
    this->mSpi = new MarketSpi(&this->UserReq,"market");
    this->tSpi = new TraderSpi(&this->UserReq,"trader");

}
ctpspi::~ctpspi(){
    delete this->mSpi;
    delete this->tSpi;
}
void ctpspi::runMRecv(){
    this->mSpi->receive();
}
void ctpspi::runTRecv(){
    this->tSpi->receive();
}
void ctpspi::runMarket(const char *addr){
    char _addr[1024];
    strcpy(_addr,addr);
    cout<<addr<<endl;
    mSpi->mdApi->RegisterFront(_addr);
    mSpi->mdApi->Init();
    mSpi->mdApi->Join();

}
void ctpspi::runTrader(const char *addr){
    char _addr[1024];
    strcpy(_addr,addr);
    cout<<addr<<endl;
    tSpi->trApi->RegisterFront(_addr);
    tSpi->trApi->Init();
    tSpi->trApi->Join();

}
