#ifndef MARKETSPI_H
#define MARKETSPI_H
#include "ThostFtdcMdApi.h"
#include "socketunixserver.h"
#include <map>
using namespace std;

class MarketSpi : public socketUnixServer,CThostFtdcMdSpi
{
public:
    MarketSpi(CThostFtdcReqUserLoginField *user, const char *path);
    CThostFtdcMdApi *mdApi;
    virtual void routeHand(const char *data);
    int getRequestID();
    virtual void OnFrontConnected();
    virtual void OnRspUserLogin(
        CThostFtdcRspUserLoginField *pRspUserLogin,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast);
    virtual void OnRspSubMarketData(
            CThostFtdcSpecificInstrumentField *pSpecificInstrument,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
            bool bIsLast) ;
    virtual void OnRspUnSubMarketData(
            CThostFtdcSpecificInstrumentField *pSpecificInstrument,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
            bool bIsLast);
    virtual void OnRtnDepthMarketData(CThostFtdcDepthMarketDataField *pDepthMarketData);


private:
    int requestID;
    CThostFtdcReqUserLoginField userReq;
    map<string , int >mapstring;

};

#endif // MARKETSPI_H
