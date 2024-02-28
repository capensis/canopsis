import { omit } from 'lodash';

import { flushPromises } from '@unit/utils/vue';
import {
  createActiveViewModule,
  createAuthModule,
  createServiceModule,
  createUserPreferenceModule,
  createWidgetModule,
  createWidgetTemplateModule,
  createInfosModule,
} from '@unit/utils/store';

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

export const submitWithExpects = async (wrapper, { fetchActiveView, hideSidebar, widgetMethod, expectData }) => {
  const widgetSettings = wrapper.vm.$children[0];

  widgetSettings.$emit('submit');

  await flushPromises();

  expect(widgetMethod).toHaveBeenCalledTimes(1);
  expect(widgetMethod).toHaveBeenLastCalledWith(
    expect.any(Object),
    expectData,
    undefined,
  );
  expect(fetchActiveView).toHaveBeenCalledTimes(1);
  expect(hideSidebar).toHaveBeenCalledTimes(1);
};
