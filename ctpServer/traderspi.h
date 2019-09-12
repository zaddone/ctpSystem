#ifndef TRADERSPI_H
#define TRADERSPI_H
#include "ThostFtdcTraderApi.h"
#include "socketunixserver.h"
#include <map>

using namespace std;
class TraderSpi:public socketUnixServer,CThostFtdcTraderSpi
{
public:
    TraderSpi(CThostFtdcReqUserLoginField *user,const char* path);
    CThostFtdcTraderApi *trApi;
    virtual void routeHand(const char *data);
    int getRequestID();
    void run(const char *addr);
    virtual void OnFrontConnected();

    virtual void OnRspUserLogin(
        CThostFtdcRspUserLoginField *pRspUserLogin,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast);
    virtual void OnRspQryInstrument(
        CThostFtdcInstrumentField *pInstrument,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast);
private:
    void swapPassword();
    TThostFtdcPasswordType pass;
    int requestID;
    CThostFtdcReqUserLoginField userReq;
    map<string , int >mapstring;
    void initMap();
    void setUserReg(
        const char * brokerID,
        const char * userID,
        const char * password,
        const char * passwordBak);
    void queryInstruments();
    bool queryIns;

};

#endif // TRADERSPI_H
