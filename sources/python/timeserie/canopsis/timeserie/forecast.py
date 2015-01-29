# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

import logging

logger = logging.getLogger('forecasting_statsMethods')


def validateSerie(aggregatedValuesByTime):
    """
    Recovery of last continued background of numerical serie
    """

    logger.debug('validate_serie')

    validity = ''

    length_AggregatedValuesByTime = len(aggregatedValuesByTime)
    aggregatedValues = [value[1] for value in aggregatedValuesByTime]

    try:
        none = aggregatedValues.index(None)
        logger.debug('there is None in aggregation list : %s', str(none))
        noneValueIndices = []

        for indice in range(length_AggregatedValuesByTime):

            if aggregatedValues[indice] is None:
                noneValueIndices.append(indice)

        logger.debug('None value number : %s' % len(noneValueIndices))

        if len(noneValueIndices) == length_AggregatedValuesByTime:
            validity = 'none'
            firstValidIndice = -1

        else:
            validity = 'by period'
            firstValidIndice = noneValueIndices[-1] + 1

    except Exception as e:
        logger.debug('Error : %s', e)
        logger.debug('No "None" in the dataset')
        validity = 'all'
        firstValidIndice = 0

    logger.debug('validity : %s' % validity)
    logger.debug('firstValidIndice : %s' % firstValidIndice)

    return {'validity': validity, 'firstValidIndice': firstValidIndice}


def detectAlerts(serie, alertList):
    """
    Detect if values in the serie correspond in alerts
    """
    logger.debug('detectAlerts')

    for alert in alertList:

        alertValue = alert[0]
        alertOperator = alert[1]
        detectedAlerts = list()

        if alertOperator == '=':

            for value in serie:

                if value[1] == alertValue:
                    detectedAlerts.append(value[0], value[1], alert)
                    break

        elif alertOperator == '<=':

            for value in serie:

                if value[1] <= alertValue:
                    detectedAlerts.append(value[0], value[1], alert)
                    break

        elif alertOperator == '>=':

            for value in serie:

                if value[1] >= alertValue:
                    detectedAlerts.append(value[0], value[1], alert)
                    break

        elif alertOperator == '<':

            for value in serie:

                if value[1] < alertValue:
                    detectedAlerts.append(value[0], value[1], alert)
                    break

        elif alertOperator == '>':

            for value in serie:

                if value[1] > alertValue:
                    detectedAlerts.append(value[0], value[1], alert)
                    break

    return detectedAlerts


def calculateRMSE(dataset1, dataset2):

    logger.debug('calculateRMSE')

    timePointNb = len(dataset1)

    if timePointNb != len(dataset2):
        logger.debug('length of datasets no compatible')
        return -1
    else:
        SCR = 0

        logger.debug('dataset1 : %s' % dataset1)
        logger.debug('dataset2 : %s' % dataset2)

        #for timePoint1,timePoint2 in dataset1,dataset2 :
        for t in range(timePointNb):
            #logger.debug( "t : %s" % t)
            timePoint1 = dataset1[t]
            timePoint2 = dataset2[t]
            #logger.debug( "tp1 : %s" % timePoint1)
            #logger.debug( "tp2 : %s" % timePoint2)
            SCR += float((timePoint1[1] - timePoint2[1]) ** 2)

        MSE = SCR / float(timePointNb)

        logger.debug('MSE : %s ' % MSE)

        RMSE = MSE ** 0.5

        logger.debug('RMSE : %s ' % RMSE)

        return RMSE


def calculateAutoCorrelation(serie):

    logger.debug('calculateAutoCorrelation')

    y0 = [0 for val in serie]
    autoCorrelationSerie = [[0,0] for val in serie]
    serieValues = [val[1] for val in serie]

    serieMean = sum(serieValues) / len(serieValues)

    for p in range(len(serie)):
        autoCorrelationSerie[p][0] = p

        for n in range(len(serie)):
            autoCorrelationSerie[p][1] += \
                (serie[n][1] - serieMean) * (serie[n - p][1] - serieMean)

    logger.debug('Auto-correlation serie : %s ' % autoCorrelationSerie)
    return autoCorrelationSerie


def calculateMobileMean(serie, degree):
    serie_size = len(serie)
    mobileAverageSerie = list()
    foot = int(degree) / 2

    if degree % 2 != 0:

        for i in range(foot, serie_size - foot):
            mobileMean = 0.0

            for f in range(- foot, foot + 1):
                mobileMean += float(serie[i + f][1])

            mobileMean /= float(degree)

            mobileAverageSerie.append([serie[i][0], mobileMean])

    else:
        for i in range( foot, serie_size - foot):

            mobileMean = float(serie[i - foot][1]) / 2.0

            for f in range(- foot + 1, foot):
                mobileMean += float(serie[i + f][1])

            mobileMean += float(serie[i + foot][1]) / 2.0
            mobileMean /= float(degree)

            mobileAverageSerie.append([serie[i][0], mobileMean])

    logger.debug('mobileAverageSerie : %s' % mobileAverageSerie)

    return mobileAverageSerie


def detectParticularPoints(serie):
    logger.debug(' detectParticularPoints ')

    summits = list()
    hollow = list()
    passes = list()

    for i in range(1, len(serie) - 1):

        if(serie[i - 1][1] < serie[i][1] and serie[i + 1][1] < serie[i][1]):
            summits.append(serie[i])

        if( serie[i-1][1] > serie[i][1] and serie[i+1][1] > serie[i][1] ):
            hollow.append(serie[i])

        if( ( serie[i][1] > 0 and serie[i+1][1] < 0 ) or
               ( serie[i][1] < 0 and serie[i+1][1] > 0 ) ):
            passes.append(serie[i])

    logger.debug(  'summits : %s' % summits )
    logger.debug( '' )
    logger.debug(  'hollow : %s' % hollow )
    logger.debug( '' )
    logger.debug(  'passes : %s' % passes )

    return {'summits' : summits,
            'hollow' : hollow,
            'passes' : passes }


def calculateSeasonality( serie ):

    autoCorrelationSerie = calculateAutoCorrelation(serie)
    autoCorrelationParticularPoints = detectParticularPoints(autoCorrelationSerie)
    passNb = len( autoCorrelationParticularPoints['passes'] )

    if passNb == 2  :
        logger.debug( 'Serie is linear and doesnt have periodicity ')
        return None

    else:
        hollowMean = 0
        hollowNb = 0
        if len( autoCorrelationParticularPoints['hollow'] ) != 0 :
            #seasonality = float( len( autoCorrelationSerie) )/ float(len( autoCorrelationParticularPoints['hollow'] ) )
            for h in xrange( len( autoCorrelationParticularPoints['hollow'])-1):
                hollowNb+=1
                hollowMean+=(autoCorrelationParticularPoints['hollow'][h+1][0]-autoCorrelationParticularPoints['hollow'][h][0])

            seasonality = int(float(hollowMean)/float(hollowNb))
        else:
            return None

        logger.debug( 'Seasonality of the serie : %s' % seasonality )

        return seasonality


def detectOutliersOrTrendChanges(serie, smoothing, maxError):

    logger.debug(  'detectOutliersOrTrendChanges')
    logger.debug(  '')
    logger.debug(  'serie : %s' % serie )
    logger.debug(  '')
    logger.debug(  'smoothing : %s' % smoothing )
    indiceRecord = {}
    indiceRecord['trendChange'] = -1
    outliers = []

    timePointNb = len(serie)

    if timePointNb != len(smoothing) :
        logger.debug(  'length of datasets no compatible' )
        return -1
    else:

        outlier = False

        for t in range(timePointNb) :
            timePoint1 = serie[t]
            timePoint2 = smoothing[t]
            RMSE = calculateRMSE( [timePoint1], [timePoint2] )

            maxValue = max(timePoint1,timePoint2)

            if maxValue != 0.0 :
                relativeError = abs(RMSE/maxValue)
            else:
                relativeError = 0.0

            logger.debug( "t : %s, relative error : %s " % (t,relativeError))

            if relativeError > maxError :

                for tbis in range(t+1,timePointNb) :
                    timePoint1 = serie[tbis]
                    timePoint2 = smoothing[tbis]
                    RMSE = calculateRMSE( [timePoint1], [timePoint2] )

                    if timePoint2[1] != 0.0 :
                        relativeError = abs(RMSE/timePoint2[1])
                    elif timePoint1[1] != 0.0 :
                        relativeError = abs(RMSE/timePoint1[1])
                    else:
                        relativeError = 0.0

                    if relativeError < maxError :
                        outliers.append(t)
                        outlier = True
                        break

                if outlier is False:
                    indiceRecord['trendChange'] = t
                    indiceRecord['outliers'] = outliers
                    return indiceRecord


    logger.debug( 'outlier list : %s' % outliers )
    indiceRecord['outliers'] = outliers
    return indiceRecord


def deleteOutliers(serie, outliers ):
    serieWithoutOutliers = []

    for indice in xrange( len(serie) ) :

        if indice not in outliers :
            serieWithoutOutliers.append(serie[indice])

    logger.debug( 'outlier list : %s' % serieWithoutOutliers )
    return serieWithoutOutliers


def repairOutliers(serie, outliers, optimizedParameters ):
    correctedSerie = list( serie )

    for indice in outliers :
        logger.debug("indice : %s" % indice )

        #Search first data in serie who isn't in outlier lsit
        for i in xrange(1, indice):
            indice_noOutlier = indice-i
            if indice_noOutlier not in outliers :
                break



        if optimizedParameters['method'] == 'h_linear' :
            forecast = calculateHoltWintersLinearMethod( correctedSerie[0:indice_noOutlier+1],
                                                         indice-indice_noOutlier,
                                                         optimizedParameters['alpha'],
                                                         optimizedParameters['beta'] )

        elif optimizedParameters['method'] == 'hw_additive' :
            forecast = calculateHoltWintersAdditiveSeasonalMethod( correctedSerie[0:indice_noOutlier+1],
                                                                   indice-indice_noOutlier,
                                                                   optimizedParameters['seasonality'],
                                                                   optimizedParameters['alpha'],
                                                                   optimizedParameters['beta'],
                                                                   optimizedParameters['gamma'] )

        elif optimizedParameters['method'] == 'hw_multiplicative' :
            forecast = calculateHoltWintersMultiplicativeSeasonalMethod( correctedSerie[0:indice_noOutlier+1],
                                                                         indice-indice_noOutlier,
                                                                         optimizedParameters['seasonality'],
                                                                         optimizedParameters['alpha'],
                                                                         optimizedParameters['beta'],
                                                                         optimizedParameters['gamma'] )

        logger.debug("old | new value : %s | %s" % (correctedSerie[indice],
                                                    forecast[-1] ) )
        correctedSerie[indice] = forecast[-1]


    return correctedSerie


def optimiseHoltWintersAlgorithm( serie, seasonality=None, method=None ):
    logger.debug('optimiseHoltWintersAlgorithm')

    alpha = [ 0.1*i for i in xrange(1,10) ]
    beta = [ 0.1*i for i in xrange(1,10) ]
    gamma = [ 0.1*i for i in xrange(1,10) ]

    if seasonality == None :
        seasonality = calculateSeasonality(serie)

    if seasonality == None :
        method = 'h_linear'

    logger.debug(" seasonality : %s" % seasonality )
    logger.debug(" method : %s" % method )

    optimizeCoeffAndMethod = { 'seasonality' : seasonality }
    optimizeCoeffAndMethod['alpha']= 0.1
    optimizeCoeffAndMethod['beta'] = 0.1

    serieLength = len(serie)
    serieTestLength = int(float(serieLength)/2.0)
    serieTest = serie[0:serieTestLength]
    duration = serieLength - serieTestLength

    if method == None :

        optimizeCoeffAndMethod['gamma'] = 0.1

        for a in alpha :

            for b in beta :

                for c in gamma :
                    addForecastingTest = calculateHoltWintersAdditiveSeasonalMethod( serieTest,
                                                                                     duration,
                                                                                     seasonality,
                                                                                     a,b,c )

                    multiForecastingTest = calculateHoltWintersMultiplicativeSeasonalMethod( serieTest,
                                                                                             duration,
                                                                                             seasonality,
                                                                                             a,b,c )

                    addError = calculateRMSE(serie,addForecastingTest[:serieLength])

                    multiError = calculateRMSE(serie,multiForecastingTest[:serieLength])

                    #logger.debug( 'addError,multiError : %s,%s ' % ( addError,multiError ) )

                    if b==0.1 and a==0.1 and c==0.1 :
                        error = addError


                    if addError > multiError and error > multiError :
                        error = multiError
                        method='hw_multiplicative'
                        optimizeCoeffAndMethod['alpha']= a
                        optimizeCoeffAndMethod['beta'] = b
                        optimizeCoeffAndMethod['gamma'] = c

                    elif multiError > addError and error > addError :
                        error = addError
                        method='hw_additive'
                        optimizeCoeffAndMethod['alpha']= a
                        optimizeCoeffAndMethod['beta'] = b
                        optimizeCoeffAndMethod['gamma'] = c

    else:

        if method == 'h_linear':

            for a in alpha :

                for b in beta :
                    #logger.debug(  'a,b : %s,%s' % (a,b) )
                    linearForecastingTest = calculateHoltWintersLinearMethod(serieTest, duration, a,b )
                    linearForecastingError = calculateRMSE(serie,linearForecastingTest[:serieLength])

                    if b==0.1 and a==0.1 :
                        error = linearForecastingError

                    #logger.debug(  'linearForecastingError : ',linearForecastingError
                    if error > linearForecastingError :
                        error = linearForecastingError
                        optimizeCoeffAndMethod['alpha']= a
                        optimizeCoeffAndMethod['beta'] = b

        elif method == 'hw_additive' :
            optimizeCoeffAndMethod['gamma'] = 0.1

            for a in alpha :

                for b in beta :

                    for c in gamma :
                        addForecastingTest = calculateHoltWintersAdditiveSeasonalMethod( serieTest,
                                                                                     duration,
                                                                                     seasonality,
                                                                                     a,b,c )

                        addError = calculateRMSE(serie,addForecastingTest[:serieLength])

                        if b==0.1 and a==0.1 and c==0.1:
                            error = addError

                        if error > addError :
                            error = addError
                            optimizeCoeffAndMethod['alpha']= a
                            optimizeCoeffAndMethod['beta'] = b
                            optimizeCoeffAndMethod['gamma'] = c


        elif method == 'hw_multiplicative' :

            optimizeCoeffAndMethod['gamma'] = 0.1

            for a in alpha :

                for b in beta :

                    for c in gamma :
                        multiForecastingTest = calculateHoltWintersMultiplicativeSeasonalMethod( serieTest,
                                                                                               duration,
                                                                                               seasonality,
                                                                                               a,b,c )

                        multiError = calculateRMSE(serie,multiForecastingTest[:serieLength])

                        if b==0.1 and a==0.1 and c==0.1 :
                            error = multiError

                        if error > multiError :
                            error = multiError
                            optimizeCoeffAndMethod['alpha']= a
                            optimizeCoeffAndMethod['beta'] = b
                            optimizeCoeffAndMethod['gamma'] = c

    optimizeCoeffAndMethod['method'] = method

    logger.debug(  'optimizeCoeffAndMethod : %s' % optimizeCoeffAndMethod )

    return optimizeCoeffAndMethod


def calculateHoltWintersLinearMethod( serie, duration, alpha=0.3, beta=0.1,callback=None ):
    logger.debug('calculateHoltWintersLinearMethod')
    #logger.debug(  "##########################################################" )
    #logger.debug(  " Initial values " )
    #logger.debug(  '%s' % serie )
    #logger.debug(  "##########################################################" )

    serieLength = len(serie)
    logger.debug( 'Length of the serie : %s' % serieLength )

    if serieLength < duration :
        logger.debug(  ' We  have a insufficient number of numerical values to calculate forecasting ' )
        return []

    else:
        forecastingDuration = serieLength + duration

        # Serie of coefficients ( trend )
        level = [0]*serieLength
        trend = [0]*serieLength
        forecastingSerie = [0]*forecastingDuration

        logger.debug(  ' forecasting duration :  %s ' % forecastingDuration )

        # Initialization
        level[0] = serie[0][1]
        trend[0] = serie[1][1] - serie[0][1]
        forecastingSerie[0] = [ serie[0][0], level[0] ]
        forecastingSerie[1] = [ serie[1][0], level[0] + trend[0] ]

        if 1 >= serieLength - duration  :
            tp = serie[0][0] + duration*( serie[-1][0]-serie[-2][0] )
            forecastingSerie[ duration ] = [ tp, level[0]+duration*trend[0] ]

        for t in xrange( 1, serieLength ):
            #logger.debug(  ' t :  %s ' % t )
            level[t] = alpha*serie[t][1] + (1-alpha)*(level[t-1]+trend[t-1])
            trend[t] = beta*(level[t]-level[t-1]) + (1-beta)*trend[t-1]

            if t+1 < serieLength :
                  forecastingSerie[ t+1 ] = [ serie[t+1][0], level[t]+trend[t] ]

            if t+1 >= serieLength - duration  :
                tp = serie[t][0] + duration*( serie[-1][0]-serie[-2][0] )
                forecastingSerie[ t+duration ] = [ tp, level[t]+duration*trend[t] ]

            #logger.debug( 'forecast serie in progress: %s ' % forecastingSerie )
        logger.debug( '  ' )
        logger.debug( ' Linear forecast serie : %s ' % forecastingSerie )

        return forecastingSerie


def calculateHoltWintersAdditiveSeasonalMethod( serie, duration, season,
                                                alpha=0.3, beta=0.3, gamma=0.3 ):

    logger.debug('calculateHoltWintersAdditiveSeasonalMethod')

    serieLength  =  len( serie )

    if serieLength < duration :
        logger.debug(  ' We have an insufficient value number to calculate forecasting ' )
        return {}

    else:
        # Serie of coefficients ( trend )
        forecastingDuration = serieLength + duration
        level = [0]*serieLength
        trend = [0]*serieLength
        seasonality =[0]*serieLength

        # forecast value list
        forecastingSerie = [0]*forecastingDuration


        # Initialization
        level[0] = sum( [ i[1] for i in serie[0:season] ]) / float(season)
        trend[0] = ( sum( [i[1] for i in serie[season:2 *season] ] ) - sum( [ i[1] for i in serie[:season] ] ))/float( season**2 )
        seasonality[0] = serie[0][1]/level[0]
        forecastingSerie[0] = [ serie[0][0], level[0] ]
        forecastingSerie[1] = [ serie[1][0], level[0] + trend[0] + seasonality[0] ]

        #logger.debug(  'level0 : %s' % level[0] )
        #logger.debug(  'trend0 : %s' % trend[0] )
        #logger.debug(  'seasonality0 : %s' % seasonality[0] )

        for t in xrange( 1, serieLength ):
            level[t] = alpha*(serie[t][1]-seasonality[t-1]) + (1-alpha)*( level[t-1]+trend[t-1])

            trend[t] = (1-beta)*trend[t-1] +beta*(level[t]-level[t-1])

            seasonality[t] = gamma*(serie[t][1]-level[t]) + (1-gamma)*seasonality[t-1]

            if( t+1 < serieLength):
                forecastingSerie[t+1] = [ serie[t+1][0], level[t] +trend[t] + seasonality[t] ]

            if t+1 >= serieLength - duration  :
                tp = serie[t][0] + (duration)*( serie[-1][0]-serie[-2][0] )

                forecastingSerie[t+duration] = [tp, level[t] + duration*trend[t] + seasonality[t] ]

                #logger.debug(  "Value for t+duration : %s" % forecastingSerie[t+duration] )

        #logger.debug( ' Additive forecast serie : %s ' % forecastingSerie )

        return forecastingSerie


def calculateHoltWintersMultiplicativeSeasonalMethod(  serie, duration,
                                                   season, alpha=0.3,
                                                   beta=0.3, gamma=0.3 ):

    logger.debug('calculateHoltWintersMultiplicativeSeasonalMethod')

    serieLength  =  len( serie )

    if serieLength < duration :
        logger.debug(  ' We  have a insufficient number of serie to calculate forecasting ' )
        return {}

    else:
        # Serie of coefficients ( trend )
        forecastingDuration = serieLength + duration
        level = [0]*serieLength
        trend = [0]*serieLength
        seasonality = [0]*serieLength

        # forecast value list
        forecastingSerie = [0]*forecastingDuration


        # Initialization
        level[0] = sum( [ i[1] for i in serie[0:season] ]) / float(season)
        trend[0] = ( sum( [i[1] for i in serie[season:2 *season] ] ) - sum( [ i[1] for i in serie[:season] ] ))/float( season**2 )
        seasonality[0] = serie[0][1]/level[0]
        forecastingSerie[0] = [ serie[0][0], level[0] ]
        forecastingSerie[1] = [ serie[1][0], level[0] + trend[0] + seasonality[0] ]
        #logger.debug(  'level0 : %s' % level[0] )
        #logger.debug(  'trend0 : %s' % trend[0] )
        #logger.debug(  'seasonality0 : %s' % seasonality[0] )

        for t in xrange( 1, serieLength ):
            level[t] = alpha*(serie[t][1]/seasonality[t-1]) + (1-alpha)*( level[t-1]+trend[t-1])
            trend[t] = (1-beta)*trend[t-1] +beta*(level[t]-level[t-1])
            seasonality[t] = gamma*(serie[t][1]/level[t]) + (1-gamma)*seasonality[t-1]

            if( t+1 < serieLength):
                forecastingSerie[t+1] = [ serie[t][0], (level[t] +trend[t])*seasonality[t] ]

            if t+1 >= serieLength - duration  :
                tp = serie[t][0] + (duration)*( serie[-1][0]-serie[-2][0] )
                forecastingSerie[t+duration] = [tp, (level[t] + duration*trend[t])*seasonality[t] ]
                #logger.debug(  "Value for t+duration : %s" % forecastingSerie[t+duration] )

        #logger.debug( ' Multiplicative forecast serie : %s ' % forecastingSerie )

        return forecastingSerie

import sys
import logging
logger = logging.getLogger('forecast')

import time

from datetime import datetime

# To delete
#from ctimeserie.threshold import Threshold


class Forecast(object):
    """
    Forecast management.
    """

    ALPHA = 'alpha'
    BETA = 'beta'
    GAMMA = 'gamma'

    DEFAULT_ALPHA = 0.99
    DEFAULT_BETA = 0.12
    DEFAULT_GAMMA = 0.80

    MAX_POINTS = 250

    NOT_LINEAR_MID_VARIABLE = 'NotLinearMidVariable'
    LINEAR_NOT_VARIABLE = 'LinearNotVariable'
    NOT_LINEAR_VARIABLE = 'NotLinearVariable'
    LINEAR_VARIABLE = 'LinearVariable'

    CURVE_TYPE = 'curve_type'

    NotLinearMidVariable = {
        CURVE_TYPE: NOT_LINEAR_MID_VARIABLE,
        ALPHA: 0.99,
        BETA: 0.12,
        GAMMA: 0.80}
    LinearNotVariable = {
        CURVE_TYPE: LINEAR_NOT_VARIABLE,
        ALPHA: 0.99,
        BETA: 0.01,
        GAMMA: 0.97}
    NotLinearVariable = {
        CURVE_TYPE: NOT_LINEAR_VARIABLE,
        ALPHA: 0.60,
        BETA: 1,
        GAMMA: 0.01}
    LinearVariable = {
        CURVE_TYPE: LINEAR_VARIABLE,
        ALPHA: 0.68,
        BETA: 0.01,
        GAMMA: 0.17}

    PARAMETERS = [
        LinearNotVariable,
        LinearVariable,
        NotLinearVariable,
        NotLinearMidVariable
    ]

    def __init__(
            self,
            timeserie,
            max_points=MAX_POINTS,
            date=None,
            duration=None,
            parameters=NotLinearMidVariable):

        super(Forecast, self).__init__()

        self.timeserie = timeserie
        self.max_points = max_points
        self.date = date
        self.duration = duration
        self.parameters = parameters

    def __repr__(self):
        message = "time_serie: %s, max_pts: %s, \
            date: %s, duration: %s, params: %s"
        result = message % (
            self.timeserie,
            self.max_points,
            self.date,
            self.duration,
            self.parameters)
        return result


    def validate_serie(y,c):
        return (len(y)>2 and len(y)%c ==0)

    def holtwinters(
            y,
            alpha=DEFAULT_ALPHA,
            beta=DEFAULT_BETA,
            gamma=DEFAULT_GAMMA,
            c=sys.maxint,
            forecast=True):
        """
        y - time series data.
        alpha(0.2) , beta(0.1), gamma(0.05): exponential smoothing coefficients
                                        for level, trend, seasonal components.
        c -  extrapolated future data points.
              4 quarterly
              7 weekly.
              12 monthly

        The length of y must be a an integer multiple (> 2) of c.
        """

        # Validate the serie
        if(not validate_serie(y,c)):
            return
        logger.debug(
            "y = %s, alpha = %s, beta = %s, gamma = %s, c = %s" %
            (y, alpha, beta, gamma, c))

        #Compute initial b and intercept using the first two complete c periods
        ylen = len(y)

        c = min(c, ylen >> 1)

        fc = float(c)

        ybar2 = sum([y[i] for i in range(c, 2 * c)]) / fc
        ybar1 = sum([y[i] for i in range(c)]) / fc
        b0 = (ybar2 - ybar1) / fc
        logger.debug("b0 = %s" % b0)

        #Compute for the level estimate a0 using b0 above.
        tbar = ((c * (c + 1)) >> 1) / fc
        logger.debug("tbar = %s" % tbar)
        a0 = ybar1 - b0 * tbar
        logger.debug("a0 = %s" % a0)

        #Compute for initial indices
        I = [y[i] / (a0 + (i + 1) * b0) for i in range(0, ylen)]
        logger.debug("Initial indices = %s" % I)

        S = [0] * (ylen + c)
        for i in range(c):
            S[i] = (I[i] + I[i + c]) / 2.0

        #Normalize so S[i] for i in [0, c)  will add to c.
        tS = c / sum([S[i] for i in range(c)])

        for i in range(c):
            S[i] *= tS
            logger.debug("S[%s]=%s" % (i, S[i]))

        # Holt - winters proper ...
        logger.debug("Use Holt Winters forecast method")
        F = [0] * (ylen + c)

        At = a0
        Bt = b0
        for i in range(ylen):
            Atm1 = At
            Btm1 = Bt
            At = alpha * y[i] / S[i] + (1.0 - alpha) * (Atm1 + Btm1)
            Bt = beta * (At - Atm1) + (1 - beta) * Btm1
            S[i + c] = gamma * y[i] / At + (1.0 - gamma) * S[i]
            F[i] = (a0 + b0 * (i + 1)) * S[i]
            logger.debug(
                "i=%s, y=%s, S=%s, Atm1=%s, Btm1=%s, \
                At=%s, Bt=%s, S[i+c]=%s, F=%s" %
                (i + 1, y[i], S[i], Atm1, Btm1, At, Bt, S[i + c], F[i]))

        if forecast:
            Forecast.holtwinters_forecast(y, c, At, Bt, F, S)

        return c, At, Bt, F

    @staticmethod
    def holtwinters_forecast(y, c, At, Bt, F, S):
        ylen = len(y)
        #Forecast for next c periods:
        for m in range(c):
            F[ylen + m] = (At + Bt * (m + 1)) * S[ylen + m]
            logger.debug("forecast: %s" % F[ylen + m])

        logger.debug("F = %s" % F)

        return F

    @staticmethod
    def forecast_best_effort(
            y,
            forecast_parameters=PARAMETERS,
            c=sys.maxint,
            calculate_forecast=True):
        """
        Identify best category for forecasting input y list.
        Returns count of forecasting points, At and Bt params and
        holtwinter list.
        This list is filled if calculate_forecast is True.
        """

        result = None

        for forecast_parameter in forecast_parameters:
            c, At, Bt, F = Forecast.holtwinters(
                y,
                forecast_parameter[Forecast.ALPHA],
                forecast_parameter[Forecast.BETA],
                forecast_parameter[Forecast.GAMMA], c)
            delta = sum(map(lambda x: abs(x[0] - x[1]), zip(y, F)))

            if result is None or result['delta'] > delta:
                result = forecast_parameter
            else:
                result['delta'] = delta
        if calculate_forecast:
            Forecast.holtwinters_forecast(y, c, At, Bt, F)

        return result, c, At, Bt, F

    def calculate_points(self, points, timewindow, timeserie):
        """
        Get forecasted points from points and forecast properties.
        """

        logger.debug('forecast: points: %s' % points)

        noneindexes = \
            [index for index, point in enumerate(points) if point[1] is None]

        logger.debug('noneindexes: %s' % noneindexes)

        # remove None values
        y = [point[1] for point in points if point[1] is not None]

        # calculate forecasted points
        if self.parameters is None:  # best effort
            forecast_parameters, c, At, Bt, F = Forecast.forecast_best_effort(
                y, Forecast.FORECAST_PARAMETERS, c)
            self.parameters = forecast_parameters

        else:
            c, At, Bt, F = Forecast.holtwinters(
                y,
                self.parameters.alpha,
                self.parameters.beta,
                self.parameters.gamma,
                self.max_points)

        # add None in F
        def insertnonevalues(F, noneindexes, index=0):
            noneindexeslen = len(noneindexes)
            for _index in range(0, noneindexeslen):
                nonindex = noneindexes[_index]
                F.insert(nonindex + index, None)

        insertnonevalues(F, noneindexes)
        insertnonevalues(F, noneindexes, len(y))

        result = [[z[0][0], z[1]] for z in zip(points, F[:len(points)])]

        date = datetime.fromtimestamp(points[-1][0])

        logger.debug('last_point: %s, date: %s' % (points[-1], date))

        for index in range(len(points), len(F)):
            date = timewindow.get_next_date(
                date,
                self.timeserie.period)
            timestamp = time.mktime(date.timetuple())
            point = [timestamp, F[index]]
            result.append(point)

        logger.debug('points: %s' % points)

        logger.debug('result: %s' % result)

        logger.debug(
            'len(points): %s, len(y): %s, len(F): %s, len(result): %s' %
            (len(points), len(y), len(F), len(result)))

        return result
