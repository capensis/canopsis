import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { createButtonStub } from '@unit/stubs/button';

import { MODALS, WIDGET_TYPES } from '@/constants';

import { alarmListChartToForm, formToAlarmListChart } from '@/helpers/entities/widget/forms/alarm';

import ChartsForm from '@/components/sidebars/chart/form/charts-form.vue';

const stubs = {
  'widget-settings-group': true,
  'field-draggable-list': true,
  'v-menu': createActivatorElementStub('v-menu'),
  'v-btn': createButtonStub('v-btn'),
  'v-list': true,
  'v-list-item': createButtonStub('v-list-item'),
};
const snapshotStubs = {
  ...stubs,

  'v-menu': true,
  'v-btn': true,
  'v-list-item': true,
};

const selectAddChartButton = wrapper => wrapper.find('.v-btn, v-btn-stub');

describe('charts-form', () => {
  const $modals = mockModals();
  const factory = generateShallowRenderer(ChartsForm, { stubs, mocks: { $modals } });
  const snapshotFactory = generateRenderer(ChartsForm, { stubs: snapshotStubs });
  const charts = Faker.datatype.array(3).map((_, index) => ({
    ...formToAlarmListChart(alarmListChartToForm()),
    key: index,
  }));

  const CHARTS_TYPES_TO_LIST_TILE_INDEXES = {
    [WIDGET_TYPES.barChart]: 0,
    [WIDGET_TYPES.lineChart]: 1,
    [WIDGET_TYPES.numbers]: 2,
  };

  test.each(Object.keys(CHARTS_TYPES_TO_LIST_TILE_INDEXES))('Show create %s chart modal', async (type) => {
    const wrapper = factory();
    const addChartButton = selectAddChartButton(wrapper);

    addChartButton.trigger('click');

    const button = wrapper.findAll('.v-list-item').at(CHARTS_TYPES_TO_LIST_TILE_INDEXES[type]);

    button.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createAlarmChart,
        config: {
          chart: { type },
          onlyExternal: true,
          title: expect.any(String),
          action: expect.any(Function),
        },
      },
    );
  });

  it('Renders `charts-form` with default props', () => {
    const wrapper = snapshotFactory();

    const dropdownContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  test('Renders `charts-form` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        charts,
      },
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });
});
