import {
  createActiveViewModule,
  createAuthModule,
  createServiceModule,
  createUserPreferenceModule,
  createWidgetModule,
  createWidgetTemplateModule,
  createInfosModule,
} from '@unit/utils/store';
import { flushPromises } from '@vue/test-utils';
import { omit } from 'lodash';

export const createSettingsMocks = () => ({
  ...createAuthModule(),
  ...createUserPreferenceModule(),
  ...createWidgetModule(),
  ...createServiceModule(),
  ...createActiveViewModule(),
  ...createWidgetTemplateModule(),
  ...createServiceModule(),
  ...createInfosModule(),
});

export const getWidgetRequestWithNewProperty = (widget, key, value) => ({
  ...omit(widget, ['_id']),

  [key]: value,
});

export const getWidgetRequestWithNewParametersProperty = (widget, key, value) => ({
  ...omit(widget, ['_id']),

  parameters: {
    ...widget.parameters,

    [key]: value,
  },
});

export const submitWithExpects = async (
  wrapper,
  { fetchActiveView, hideSidebar, widgetMethod, expectData },
) => {
  const widgetSettingsElement = wrapper.find('widget-settings-stub');

  if (widgetSettingsElement.exists()) {
    widgetSettingsElement.triggerCustomEvent('submit');
  } else {
    const widgetSettingsWrapper = wrapper.findComponent({ name: 'widget-settings' });
    const widgetSettings = widgetSettingsWrapper.exists()
      ? widgetSettingsWrapper.vm
      : wrapper.vm.$children[0];

    widgetSettings.$emit('submit');
  }

  vi.runOnlyPendingTimers();
  await flushPromises();

  expect(widgetMethod).toHaveBeenCalledTimes(1);
  expect(widgetMethod).toHaveBeenLastCalledWith(expect.any(Object), expectData);
  expect(fetchActiveView).toHaveBeenCalledTimes(1);
  expect(hideSidebar).toHaveBeenCalledTimes(1);
};
