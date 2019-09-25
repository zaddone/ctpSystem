#ifndef TRADERSPI_H
#define TRADERSPI_H
#include "ThostFtdcTraderApi.h"
#include "socketunixserver.h"
#include <map>

using namespace std;
class TraderSpi:public socketUnixServer,CThostFtdcTraderSpi
{
public:
    TraderSpi(const char* path);
    //TraderSpi(CThostFtdcReqUserLoginField *user,const char* path);
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

    virtual void OnRspSettlementInfoConfirm(
            CThostFtdcSettlementInfoConfirmField *pSettlementInfoConfirm,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
            bool bIsLast) ;
    virtual void OnRspQrySettlementInfo(
            CThostFtdcSettlementInfoField *pSettlementInfo,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
            bool bIsLast);
    virtual void OnRspQrySettlementInfoConfirm(
            CThostFtdcSettlementInfoConfirmField *pSettlementInfoConfirm,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
            bool bIsLast);

    virtual void OnRspOrderInsert(
            CThostFtdcInputOrderField *pInputOrder,
            CThostFtdcRspInfoField *pRspInfo, int nRequestID, bool bIsLast);
    virtual void OnRtnOrder(CThostFtdcOrderField *pOrder);
    virtual void OnRtnTrade(CThostFtdcTradeField *pTrade);
    virtual void OnErrRtnOrderInsert(CThostFtdcInputOrderField *pInputOrder, CThostFtdcRspInfoField *pRspInfo);
    virtual void OnRspQryTradingAccount(CThostFtdcTradingAccountField *pTradingAccount,
                                        CThostFtdcRspInfoField *pRspInfo,
                                        int nRequestID,
                                        bool bIsLast) ;

    virtual void OnRspQryInvestorPositionDetail(
            CThostFtdcInvestorPositionDetailField *pInvestorPositionDetail,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
            bool bIsLast);
    virtual void OnRspQryInvestorPosition(
            CThostFtdcInvestorPositionField *pInvestorPosition,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
            bool bIsLast);

private:
    void swapPassword();
    TThostFtdcPasswordType pass;
    int requestID;
    CThostFtdcReqUserLoginField userReq;
    map<string , int >mapstring;
    map<string, CThostFtdcInstrumentField >mapInstrument;
    void initMap();
    void Join();
    void setUserReg(
        const char * brokerID,
        const char * userID,
        const char * password,
        const char * passwordBak);
    void reqUserLogin();
    void reqInstruments();
    void reqSettlementInfoConfirm();
    void reqQrySettlementInfo();
    void reqQrySettlementInfoConfirm();
    void reqTradingAccount();
    void reqInvestorPosition(const char *ins);
    void reqInvestorPositionDetail(const char * ins);

    void stop();
    void sendOrderOpen(const char * ins, const char *dir, const double price=0);
    void sendOrderClose(const char * ins);
    //void investorPosition(const char * ins);
    const char * path;
    bool queryIns;

};

#endif // TRADERSPI_H
