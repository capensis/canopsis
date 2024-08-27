import { omit } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockSidebar } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createInputStub } from '@unit/stubs/input';
import {
  createSettingsMocks,
  getWidgetRequestWithNewProperty,
  getWidgetRequestWithNewParametersProperty,
  submitWithExpects,
} from '@unit/utils/settings';

import { SIDE_BARS, WIDGET_TYPES } from '@/constants';

import ClickOutside from '@/services/click-outside';

import { widgetToForm, formToWidget, getEmptyWidgetByType } from '@/helpers/entities/widget/form';

import TextSettings from '@/components/sidebars/text/text.vue';

const stubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': createInputStub('field-title'),
  'field-text-editor': createInputStub('field-text-editor'),
  'field-stats-selector': createInputStub('field-stats-selector'),
  'field-date-interval': createInputStub('field-date-interval'),
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-text-editor': true,
  'field-stats-selector': true,
  'field-date-interval': true,
};

const generateDefaultTextWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.text })),

  _id: Faker.datatype.string(),
});

const selectFieldTitle = wrapper => wrapper.find('input.field-title');
const selectFieldTextEditor = wrapper => wrapper.find('input.field-text-editor');

describe('text', () => {
  const $sidebar = mockSidebar();

  const {
    createWidget,
    updateWidget,
    fetchActiveView,
    fetchUserPreference,
    activeViewModule,
    widgetModule,
    authModule,
    userPreferenceModule,
  } = createSettingsMocks();

  const widget = {
    ...generateDefaultTextWidget(),

    tab: Faker.datatype.string(),
  };

  const sidebar = {
    name: SIDE_BARS.textSettings,
    config: {
      widget,
    },
    hidden: false,
  };

  const store = createMockedStoreModules([
    activeViewModule,
    widgetModule,
    authModule,
    userPreferenceModule,
  ]);

  const factory = generateShallowRenderer(TextSettings, { stubs,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(TextSettings, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  beforeAll(() => jest.useFakeTimers());

  afterEach(() => {
    createWidget.mockReset();
    updateWidget.mockReset();
    fetchActiveView.mockReset();
    fetchUserPreference.mockReset();
  });

  it('Create widget with default parameters', async () => {
    const localWidget = getEmptyWidgetByType(WIDGET_TYPES.text);

    localWidget.tab = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            widget: localWidget,
          },
        },
      },
      mocks: {
        $sidebar,
      },
    });

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: createWidget,
      expectData: {
        data: {
          ...formToWidget(widgetToForm(localWidget)),

          tab: localWidget.tab,
        },
      },
    });
  });

  it('Duplicate widget without changes', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            widget,
            duplicate: true,
          },
        },
      },
      mocks: {
        $sidebar,
      },
    });

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: createWidget,
      expectData: {
        data: omit(widget, ['_id']),
      },
    });
  });

  it('Title changed after trigger field title', async () => {
    const newTitle = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldTitle = selectFieldTitle(wrapper);

    fieldTitle.setValue(newTitle);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(widget, 'title', newTitle),
      },
    });
  });

  it('Template changed after trigger field text editor', async () => {
    const newTemplate = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldTextEditor = selectFieldTextEditor(wrapper);

    fieldTextEditor.setValue(newTemplate);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'template', newTemplate),
      },
    });
  });

  it('Renders `text` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `text` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            ...sidebar.config,

            widget: {
              _id: '_id123',
              type: WIDGET_TYPES.text,
              title: 'Text widget',
              parameters: {
                template: '<p>1</p><p>2</p><p>3<br></p>',
              },
            },
          },
        },
      },
      mocks: {
        $sidebar,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
