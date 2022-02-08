import { omit } from 'lodash';
import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createButtonStub } from '@unit/stubs/button';
import { createInputStub } from '@unit/stubs/input';

import { CANOPSIS_EDITION, QUICK_RANGES, TIME_UNITS, WIDGET_TYPES } from '@/constants';
import ClickOutside from '@/services/click-outside';
import { widgetToForm, formToWidget } from '@/helpers/forms/widgets/common';

import TextSettings from '@/components/sidebars/settings/text.vue';

const localVue = createVueInstance();

const stubs = {
  'field-title': createInputStub('field-title'),
  'field-filter-editor': createInputStub('field-filter-editor'),
  'field-text-editor': createInputStub('field-text-editor'),
  'field-stats-selector': createInputStub('field-stats-selector'),
  'field-date-interval': createInputStub('field-date-interval'),
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'field-title': true,
  'field-filter-editor': true,
  'field-text-editor': true,
  'field-stats-selector': true,
  'field-date-interval': true,
};

const factory = (options = {}) => shallowMount(TextSettings, {
  localVue,
  stubs,

  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(TextSettings, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectSubmitButton = wrapper => wrapper.find('button.v-btn');
const selectFieldTitle = wrapper => wrapper.find('input.field-title');
const selectFieldTextEditor = wrapper => wrapper.find('input.field-text-editor');
const selectFieldFilterEditor = wrapper => wrapper.find('input.field-filter-editor');
const selectFieldStatsSelector = wrapper => wrapper.find('input.field-stats-selector');
const selectFieldDateInterval = wrapper => wrapper.find('input.field-date-interval');

describe('text', () => {
  const generateDefaultTextWidget = () => ({
    ...formToWidget(widgetToForm({ type: WIDGET_TYPES.text })),

    _id: Faker.datatype.string(),
  });

  const userPreferences = {
    content: Faker.helpers.createTransaction(),
  };
  const widget = {
    ...generateDefaultTextWidget(),

    tab: Faker.datatype.string(),
  };
  const view = {
    enabled: true,
    title: 'Text widgets',
    description: 'Text widgets',
    tabs: [
      {
        widgets: [widget],
      },
    ],
    tags: ['text'],
    periodic_refresh: {
      value: 1,
      unit: 's',
      enabled: false,
    },
    author: 'root',
    group: {
      _id: 'text-widget-group',
    },
  };
  const updateUserPreference = jest.fn();
  const updateView = jest.fn();
  const updateQuery = jest.fn();
  const hideSideBar = jest.fn();

  const store = createMockedStoreModules([
    {
      name: 'sideBar',
      actions: {
        hide: hideSideBar,
      },
    },
    {
      name: 'info',
      getters: { edition: CANOPSIS_EDITION.cat },
    },
    {
      name: 'query',
      getters: { getQueryById: () => () => ({}) },
      actions: {
        update: updateQuery,
      },
    },
    {
      name: 'view',
      getters: {
        item: view,
      },
      actions: {
        update: updateView,
      },
    },
    {
      name: 'userPreference',
      getters: {
        getItemByWidget: () => () => userPreferences,
      },
      actions: {
        update: updateUserPreference,
      },
    },
  ]);

  const submitWithExpects = async (wrapper, { widgetMethod, expectData }) => {
    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(widgetMethod).toHaveBeenCalledTimes(1);
    expect(widgetMethod).toHaveBeenLastCalledWith(
      expect.any(Object),
      expectData,
      undefined,
    );
    expect(fetchActiveView).toHaveBeenCalledTimes(1);
    expect($sidebar.hide).toHaveBeenCalledTimes(1);
  };

  afterEach(() => {
    updateUserPreference.mockReset();
    updateView.mockReset();
    updateQuery.mockReset();
    hideSideBar.mockReset();
  });

  it('Title changed after trigger field title', async () => {
    const newTitle = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        sidebar: {
          config: {
            widget,
          },
        },
      },
    });

    const fieldTitle = selectFieldTitle(wrapper);
    const submitButton = selectSubmitButton(wrapper);

    fieldTitle.setValue(newTitle);

    submitButton.trigger('click');

    await flushPromises();

    expect(updateView).toHaveBeenCalledTimes(1);
    expect(updateView).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: view._id,
        data: {
          ...omit(view, ['_id', 'tabs']),
          tabs: [{
            widgets: [{
              ...widget,
              title: newTitle,
            }],
          }],
          group: view.group._id,
        },
      },
      undefined,
    );

    expect(updateQuery).toHaveBeenCalledTimes(1);
    expect(updateQuery).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          dateInterval: widget.parameters.dateInterval,
          mfilter: widget.parameters.mfilter,
          stats: widget.parameters.stats,
          template: widget.parameters.template,
        },
      },
      undefined,
    );

    expect(hideSideBar).toHaveBeenCalledTimes(1);
  });

  it('Filter changed after trigger field filter editor', async () => {
    const $clickOutside = new ClickOutside();
    const filter = {
      filter: {
        $and: [{ impact_state: { $gt: 2 } }],
      },
      title: 'Filter title',
    };

    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
      parentComponent: {
        provide: {
          $clickOutside,
        },
      },
    });

    const fieldFilterEditor = selectFieldFilterEditor(wrapper);
    const submitButton = selectSubmitButton(wrapper);

    fieldFilterEditor.vm.$emit('input', filter);

    submitButton.trigger('click');

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledTimes(1);
    expect(updateUserPreference).toHaveBeenLastCalledWith(
      expect.any(Object),
      { data: userPreferences },
      undefined,
    );

    expect(updateView).toHaveBeenCalledTimes(1);
    expect(updateView).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: view._id,
        data: {
          ...omit(view, ['_id', 'tabs']),
          tabs: [{
            widgets: [{
              ...widget,
              parameters: {
                ...widget.parameters,
                mfilter: filter,
              },
            }],
          }],
          group: view.group._id,
        },
      },
      undefined,
    );

    expect(updateQuery).toHaveBeenCalledTimes(1);
    expect(updateQuery).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          dateInterval: widget.parameters.dateInterval,
          mfilter: filter,
          stats: widget.parameters.stats,
          template: widget.parameters.template,
        },
      },
      undefined,
    );

    expect(hideSideBar).toHaveBeenCalledTimes(1);
  });

  it('Template changed after trigger field text editor', async () => {
    const $clickOutside = new ClickOutside();
    const newTemplate = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
      parentComponent: {
        provide: {
          $clickOutside,
        },
      },
    });

    const fieldTextEditor = selectFieldTextEditor(wrapper);
    const submitButton = selectSubmitButton(wrapper);

    fieldTextEditor.setValue(newTemplate);

    submitButton.trigger('click');

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledTimes(1);
    expect(updateUserPreference).toHaveBeenLastCalledWith(
      expect.any(Object),
      { data: userPreferences },
      undefined,
    );

    expect(updateView).toHaveBeenCalledTimes(1);
    expect(updateView).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: view._id,
        data: {
          ...omit(view, ['_id', 'tabs']),
          tabs: [{
            widgets: [{
              ...widget,
              parameters: {
                ...widget.parameters,
                template: newTemplate,
              },
            }],
          }],
          group: view.group._id,
        },
      },
      undefined,
    );

    expect(updateQuery).toHaveBeenCalledTimes(1);
    expect(updateQuery).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          dateInterval: widget.parameters.dateInterval,
          mfilter: widget.parameters.mfilter,
          stats: widget.parameters.stats,
          template: newTemplate,
        },
      },
      undefined,
    );

    expect(hideSideBar).toHaveBeenCalledTimes(1);
  });

  it('Stats changed after trigger field stats selector', async () => {
    const $clickOutside = new ClickOutside();
    const name = Faker.datatype.string();
    const availableStats = {
      stat: {
        value: Faker.datatype.string(),
        options: [Faker.datatype.string(), Faker.datatype.string()],
      },
      trend: Faker.datatype.boolean(),
      parameters: {
        recursive: Faker.datatype.boolean(),
        states: [0],
      },
    };
    const stats = {
      [name]: availableStats,
    };

    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
      parentComponent: {
        provide: {
          $clickOutside,
        },
      },
    });

    const fieldStatsSelector = selectFieldStatsSelector(wrapper);
    const submitButton = selectSubmitButton(wrapper);

    fieldStatsSelector.vm.$emit('input', stats);

    submitButton.trigger('click');

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledTimes(1);
    expect(updateUserPreference).toHaveBeenLastCalledWith(
      expect.any(Object),
      { data: userPreferences },
      undefined,
    );

    expect(updateView).toHaveBeenCalledTimes(1);
    expect(updateView).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: view._id,
        data: {
          ...omit(view, ['_id', 'tabs']),
          tabs: [{
            widgets: [{
              ...widget,
              parameters: {
                ...widget.parameters,
                stats,
              },
            }],
          }],
          group: view.group._id,
        },
      },
      undefined,
    );

    expect(updateQuery).toHaveBeenCalledTimes(1);
    expect(updateQuery).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          stats: {
            [name]: {
              ...availableStats,
              stat: availableStats.stat.value,
            },
          },
          dateInterval: widget.parameters.dateInterval,
          mfilter: widget.parameters.mfilter,
          template: widget.parameters.template,
        },
      },
      undefined,
    );

    expect(hideSideBar).toHaveBeenCalledTimes(1);
  });

  it('Interval changed after trigger field date interval', async () => {
    const $clickOutside = new ClickOutside();
    const dateInterval = {
      periodUnit: TIME_UNITS.week,
      periodValue: 2,
      tstart: QUICK_RANGES.last30Days.start,
      tstop: QUICK_RANGES.last30Days.stop,
    };

    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
      parentComponent: {
        provide: {
          $clickOutside,
        },
      },
    });

    const fieldDateInterval = selectFieldDateInterval(wrapper);
    const submitButton = selectSubmitButton(wrapper);

    fieldDateInterval.vm.$emit('input', dateInterval);

    submitButton.trigger('click');

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledTimes(1);
    expect(updateUserPreference).toHaveBeenLastCalledWith(
      expect.any(Object),
      { data: userPreferences },
      undefined,
    );

    expect(updateView).toHaveBeenCalledTimes(1);
    expect(updateView).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: view._id,
        data: {
          ...omit(view, ['_id', 'tabs']),
          tabs: [{
            widgets: [{
              ...widget,
              parameters: {
                ...widget.parameters,
                dateInterval,
              },
            }],
          }],
          group: view.group._id,
        },
      },
      undefined,
    );

    expect(updateQuery).toHaveBeenCalledTimes(1);
    expect(updateQuery).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          dateInterval,
          stats: widget.parameters.stats,
          mfilter: widget.parameters.mfilter,
          template: widget.parameters.template,
        },
      },
      undefined,
    );

    expect(hideSideBar).toHaveBeenCalledTimes(1);
  });

  it('Renders `text` widget settings with default props', async () => {
    const $clickOutside = new ClickOutside();

    const wrapper = snapshotFactory({
      propsData: {
        config: {
          widget: generateDefaultTextWidgetForm(),
        },
      },
      store: createMockedStoreModules([
        {
          name: 'info',
          getters: { edition: CANOPSIS_EDITION.cat },
        },
      ]),
      parentComponent: {
        provide: {
          $clickOutside,
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `text` widget settings with core version', async () => {
    const $clickOutside = new ClickOutside();

    const wrapper = snapshotFactory({
      propsData: {
        config: {
          widget: generateDefaultTextWidgetForm(),
        },
      },
      store: createMockedStoreModules([
        {
          name: 'info',
          getters: { edition: CANOPSIS_EDITION.core },
        },
      ]),
      parentComponent: {
        provide: {
          $clickOutside,
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `text` widget settings with custom props', async () => {
    const $clickOutside = new ClickOutside();

    const wrapper = snapshotFactory({
      propsData: {
        config: {
          widget: {
            _id: '_id123',
            type: WIDGET_TYPES.text,
            title: 'Text widget',
            parameters: {
              dateInterval: {
                periodValue: 2,
                periodUnit: TIME_UNITS.week,
                tstart: QUICK_RANGES.last30Days.start,
                tstop: QUICK_RANGES.last30Days.stop,
              },
              mfilter: {
                filter: {
                  $and: [{ impact_state: { $gt: 2 } }],
                },
                title: 'Filter title',
              },
              stats: {
                Test: {
                  stat: {
                    value: 'alarms_created',
                    options: ['recursive', 'states', 'authors'],
                  },
                  trend: false,
                  parameters: { recursive: true, states: [0] },
                },
              },
              template: '<p>1</p><p>2</p><p>3<br></p>',
            },
          },
        },
      },
      store: createMockedStoreModules([
        {
          name: 'info',
          getters: { edition: CANOPSIS_EDITION.cat },
        },
      ]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
