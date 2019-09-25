#ifndef CTPSPI_H
#define CTPSPI_H
#include "marketspi.h"
#include "traderspi.h"


class ctpspi
{
public:
    ctpspi();
    //ctpspi(const char *BrokerID, const char *UserID, const char * Password);
    ~ctpspi();
    void runTrader(const char *addr);
    void runMarket(const char *addr);
    void runMRecv();
    void runTRecv();
    void getConfigM();
    void getConfigT();

private:
    CThostFtdcReqUserLoginField UserReq;
    MarketSpi * mSpi;
    TraderSpi * tSpi;

};

#endif // CTPSPI_H
