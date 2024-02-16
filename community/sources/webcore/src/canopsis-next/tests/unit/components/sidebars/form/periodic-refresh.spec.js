import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { TIME_UNITS } from '@/constants';

import PeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';

const stubs = {
  'widget-settings-item': true,
  'periodic-refresh-field': true,
  'live-watching-field': true,
};

const selectPeriodicRefreshField = wrapper => wrapper.find('periodic-refresh-field-stub');
const selectLiveWatchingField = wrapper => wrapper.find('live-watching-field-stub');

describe('periodic-refresh', () => {
  const factory = generateShallowRenderer(PeriodicRefresh, { stubs });
  const snapshotFactory = generateRenderer(PeriodicRefresh, {
    stubs,
    parentComponent: {
      provide: {
        list: {
          register: jest.fn(),
          unregister: jest.fn(),
        },
        listClick: jest.fn(),
      },
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  it('Unit as seconds settled created, if unit doesn\'t exist', () => {
    const form = {
      periodic_refresh: {
        value: 1,
      },
    };

    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      ...form,
      periodic_refresh: {
        ...form.periodic_refresh,

        unit: TIME_UNITS.second,
      },
    });
  });

  it('Value changed after trigger periodic refresh field', () => {
    const wrapper = factory({
      propsData: {
        form: {
          periodic_refresh: {
            value: 1,
            unit: TIME_UNITS.day,
          },
        },
      },
    });

    const newValue = {
      value: 2,
      unit: TIME_UNITS.week,
    };

    selectPeriodicRefreshField(wrapper).vm.$emit('input', newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      periodic_refresh: newValue,
    });
  });

  it('Live watching triggers input on changes', () => {
    const form = {
      periodic_refresh: {
        value: 1,
        unit: TIME_UNITS.second,
      },
    };

    const wrapper = factory({
      propsData: {
        form,
        withLiveWatching: true,
      },
    });

    const newLiveWatching = true;

    selectLiveWatchingField(wrapper).vm.$emit('input', newLiveWatching);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      ...form,
      liveWatching: newLiveWatching,
    });
  });

  it('Renders `periodic-refresh` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `periodic-refresh` with with life watching', () => {
    const wrapper = snapshotFactory({
      propsData: {
        withLifeWatching: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `periodic-refresh` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        name: 'customName',
        form: {
          periodic_refresh: {
            value: 1,
            unit: TIME_UNITS.minute,
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `periodic-refresh` with errors', async () => {
    const name = 'custom-name';

    const wrapper = snapshotFactory({
      propsData: {
        form: {
          periodic_refresh: {
            value: 1,
            unit: TIME_UNITS.minute,
          },
        },
        name,
      },
    });

    const validator = wrapper.getValidator();

    const periodicRefreshField = selectPeriodicRefreshField(wrapper);

    validator.attach({
      name,
      rules: 'required:true',
      getter: () => true,
      context: () => periodicRefreshField.vm,
      vm: periodicRefreshField.vm,
    });

    validator.errors.add({
      field: name,
      msg: 'error-message',
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
