import { storeToRefs } from 'pinia';

import { useWeatherServiceWidgetStore } from '@/stores/modules/weather-service-widget';

/**
 * Weather service widget store facade.
 *
 * @returns {Object}
 */
export const useWeatherServiceWidget = () => {
  const weatherServiceWidgetStore = useWeatherServiceWidgetStore();
  const { widgets, getListByWidgetId, getPendingByWidgetId, getErrorByWidgetId } =
    storeToRefs(weatherServiceWidgetStore);

  return {
    weatherServicesByWidget: widgets,
    getWeatherServicesListByWidgetId: getListByWidgetId,
    getWeatherServicesPendingByWidgetId: getPendingByWidgetId,
    getWeatherServicesErrorByWidgetId: getErrorByWidgetId,
    fetchWeatherServicesList: weatherServiceWidgetStore.fetchList,
  };
};
