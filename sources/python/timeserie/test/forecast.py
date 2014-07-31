"""
import unittest
import logging
import math as m


logging.basicConfig(
	level=logging.INFO, format='%(name)s %(levelname)s %(message)s')

from forecastingMethods import validateSerie, calculateRMSE, \
	calculateMobileMean, detectParticularPoints, calculateSeasonality, \
	detectOutliers


class forecastingMethodsToTest(unittest.TestCase):

	def setUp(self):
		pass

	def test_validate_serie(self):

		serie_all_validity = []
		serie_validity_by_period = []
		serie_none_validity = []

		for i in range(10):
			serie_all_validity.append([i * 60, i])
			serie_none_validity.append([i * 60, None])

			if i == 4:
				serie_validity_by_period.append([i * 60, None])
			else:
				serie_validity_by_period.append([i * 60, i])

		validity = validateSerie(serie_all_validity)

		assert validity['validity'] == 'all', 'all : good validity no detected'
		assert validity['firstValidIndice'] == 0, ' firstValidIndice 0 no detected'

		validity = validateSerie(serie_none_validity)

		assert validity['validity'] == 'none', 'none : good validity no detected'
		assert validity['firstValidIndice'] == -1, ' firstValidIndice -1 no detected'

		validity = validateSerie(serie_validity_by_period)

		assert validity['validity'] == 'by period', 'by period : good validity no detected'
		assert validity['firstValidIndice'] == 5, ' firstValidIndice ( 5  ) no detected'

	def test_calculate_RMSE(self):

		dataset1 = []
		dataset2 = []
		test_value = 0
		for i in range(10):
			test_value += i ** 2
			dataset1.append([60 * i, i])
			dataset2.append([60 * i, 2 * i])

		RMSE = calculateRMSE(dataset1, dataset2)
		assert RMSE == 28.5 ** 0.5, 'meanSqrstError : incorrected calculing '

	def test_calculate_mobile_mean(self):
		serie = []

		for i in range(10):
			serie.append([i * 60, i])

		mobileMean2M = calculateMobileMean(serie, 2)
		assert mobileMean2M == serie[1:9], '2M mobile mean is falser'

		mobileMean2M1 = calculateMobileMean(serie, 3)
		assert mobileMean2M1 == serie[1:9], '2M1 mobile mean is false'

	def test_detect_particular_points(self):

		serie = [[0, 0], [60, 2], [120, 4], [180, 3], [240, -1],
			[300, -5], [360, 4], [420, 3], [480, 5], [540, 7]]

		particularPts = detectParticularPoints(serie)

		assert particularPts['passes'] == [[180, 3], [300, -5]], "we haven't find good passes"
		assert particularPts['summits'] == [[120, 4], [360, 4]], "we haven't find good summits"
		assert particularPts['hollow'] == [[300, -5], [420, 3]], "we haven't find good hollow"

	def test_calculate_seasonality(self):
		serie_linear = []
		serie_periodic = []

		for i in range(100):
			serie_linear.append([60 * i, i])
			serie_periodic.append([60 * i, m.cos(i)])

		linearSeasonality = calculateSeasonality(serie_linear)
		assert linearSeasonality is None, 'Linearity no detected'

		periodicSeasonality = calculateSeasonality(serie_periodic)
		assert periodicSeasonality == int(2 * m.pi), 'No good value for Seasonality '

	def test_detect_outliers(self):

		serie = []
		smoothing = []
		for i in range(100):

			if i in (3, 8):
				serie.append([60 * i, i + 2])
			else:
				serie.append([60 * i, i])

			smoothing.append([60 * i, i])

		outliers = detectOutliers(serie, smoothing, 0.15)
		assert outliers == [3, 8], 'Outliers no detected'

if __name__ == "__main__":

	unittest.main(verbosity=2)
"""