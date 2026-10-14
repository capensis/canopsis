const feature = vi.hoisted(() => {
  const componentName = 'feature-widget';
  const modalName = 'feature-modal';
  const sidebarName = 'feature-widget-settings';
  const widgetType = 'FeatureWidget';
  const component = { name: 'FeatureWidget' };
  const componentLoader = vi.fn(async () => ({ default: component }));
  const modalLoader = vi.fn(async () => ({ default: { name: 'FeatureModal' } }));
  const sidebarLoader = vi.fn(async () => ({ default: { name: sidebarName } }));
  const preparer = vi.fn(function () {
    return { ...this.widget, preparedByFeature: true };
  });
  const parametersToForm = vi.fn(parameters => ({ ...parameters, preparedForForm: true }));
  const formToParameters = vi.fn(parameters => ({ ...parameters, preparedForApi: true }));
  const values = {
    i18n: { en: { featureWidget: { title: 'Feature widget' } } },
    'constants.WIDGET_TYPES': { featureWidget: widgetType },
    'constants.COMPONENTS_BY_WIDGET_TYPES': { [widgetType]: componentName },
    'components.modals.components': { [modalName]: modalLoader },
    'components.sidebars.components': { [sidebarName]: sidebarLoader },
    'components.widgetWrapper.components': { [componentName]: componentLoader },
    'components.widgetWrapper.computed.preparedWidget': preparer,
    'helpers.forms.widgets.widgetParametersToForm.widgetsMap': {
      [widgetType]: parametersToForm,
    },
    'helpers.forms.widgets.formToWidgetParameters.widgetsMap': {
      [widgetType]: formToParameters,
    },
  };

  return {
    componentLoader,
    componentName,
    formToParameters,
    modalLoader,
    modalName,
    parametersToForm,
    preparer,
    sidebarLoader,
    sidebarName,
    values,
    widgetType,
  };
});

vi.mock('@/services/features', () => ({
  featuresService: {
    get: (key, defaultValue) => feature.values[key] ?? defaultValue,
  },
}));

import { TOP_LEVEL_WIDGET_TYPES } from '@/constants';

import { i18n } from '@/i18n';

import { formToWidgetParameters, widgetParametersToForm } from '@/helpers/entities/widget/form';

import SidebarBase from '@/components/layout/sidebar/sidebar-base.vue';
import RegisteredModals from '@/components/modals';
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

describe('widget-wrapper feature integration', () => {
  it('registers and resolves a widget component supplied by a feature', () => {
    const component = WidgetWrapper.components[feature.componentName];

    expect(component).not.toBe(feature.componentLoader);
    expect(WidgetWrapper.computed.component.call({ widget: { type: feature.widgetType } })).toBe(
      component,
    );
  });

  it('prepares a custom widget with the feature preparer', () => {
    const widget = { type: feature.widgetType };
    const context = { widget };

    expect(WidgetWrapper.computed.preparedWidget.call(context)).toEqual({
      ...widget,
      preparedByFeature: true,
    });
    expect(feature.preparer).toHaveBeenCalledOnce();
    expect(feature.preparer.mock.contexts[0]).toBe(context);
  });

  it('makes a feature widget available for creation', () => {
    expect(TOP_LEVEL_WIDGET_TYPES).toContain(feature.widgetType);
    expect(i18n.global.t('featureWidget.title', {}, { locale: 'en' })).toBe('Feature widget');
  });

  it('uses feature converters for widget settings forms', () => {
    const parameters = { value: 'feature-value' };

    expect(widgetParametersToForm({ type: feature.widgetType, parameters })).toMatchObject({
      ...parameters,
      preparedForForm: true,
    });
    expect(formToWidgetParameters({ type: feature.widgetType, parameters })).toMatchObject({
      ...parameters,
      preparedForApi: true,
    });
    expect(feature.parametersToForm).toHaveBeenCalledWith(parameters);
    expect(feature.formToParameters).toHaveBeenCalledWith(parameters);
  });

  it('registers feature sidebar and modal components', () => {
    const sidebarComponent = SidebarBase.components[feature.sidebarName];
    const modalComponent = RegisteredModals[feature.modalName];

    expect(sidebarComponent).not.toBe(feature.sidebarLoader);
    expect(modalComponent).not.toBe(feature.modalLoader);
    expect(SidebarBase.computed.component.call({ sidebar: { name: feature.sidebarName } })).toBe(
      sidebarComponent,
    );
  });
});
