#ifndef TRADERSPI_H
#define TRADERSPI_H
#include "ThostFtdcTraderApi.h"
#include "socketunixserver.h"

class TraderSpi:public socketUnixServer,CThostFtdcTraderSpi
{
public:
    TraderSpi(CThostFtdcReqUserLoginField *user,const char* path);
    CThostFtdcTraderApi *trApi;
    virtual void routeHand(const char *data);
    int getRequestID();
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
    int requestID;
    CThostFtdcReqUserLoginField userReq;
    void queryInstruments();
    bool queryIns;
};

#endif // TRADERSPI_H
