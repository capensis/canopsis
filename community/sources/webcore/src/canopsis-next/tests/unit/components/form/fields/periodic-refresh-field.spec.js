import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';
import { TIME_UNITS } from '@/constants';

const localVue = createVueInstance();

const stubs = {
  'c-enabled-field': {
    props: ['value'],
    inject: ['$validator'],
    template: `
      <input
        v-validate="'required'"
        class="c-enabled-field"
        name="enabled"
        :value="value"
      />
    `,
  },
  'c-duration-field': {
    props: ['name', 'value'],
    inject: ['$validator'],
    template: `
      <input
        v-validate="'required'"
        class="c-duration-field"
        :name="\`\${name}.value\`"
        :value="value"
      />
    `,
  },
};

const snapshotStubs = {
  'c-enabled-field': true,
  'c-duration-field': true,
};

const factory = (options = {}) => shallowMount(PeriodicRefreshField, {
  localVue,
  stubs,
  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(PeriodicRefreshField, {
  localVue,
  stubs: snapshotStubs,
  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const selectEnabledField = wrapper => wrapper.find('.c-enabled-field');
const selectDurationField = wrapper => wrapper.find('.c-duration-field');

describe('periodic-refresh-field', () => {
  it('Enabled changed after trigger the enabled field', () => {
    const periodicRefresh = {
      enabled: true,
      unit: TIME_UNITS.second,
      value: 2,
    };
    const wrapper = factory({
      propsData: {
        periodicRefresh,
      },
    });

    const enabledField = selectEnabledField(wrapper);

    enabledField.vm.$emit('input', false);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      ...periodicRefresh,
      enabled: false,
    });
  });

  it('Duration changed after trigger the duration field', () => {
    const periodicRefresh = {
      enabled: true,
      unit: TIME_UNITS.second,
      value: 2,
    };
    const wrapper = factory({
      propsData: {
        periodicRefresh,
      },
    });

    const durationField = selectDurationField(wrapper);

    const newDuration = {
      enabled: true,
      unit: TIME_UNITS.week,
      value: 5,
    };

    durationField.vm.$emit('input', newDuration);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual(newDuration);
  });

  it('Renders `periodic-refresh-field` with default props', () => {
    const periodicRefresh = {
      enabled: true,
      unit: TIME_UNITS.second,
      value: 2,
    };
    const wrapper = snapshotFactory({
      propsData: {
        periodicRefresh,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `periodic-refresh-field` with custom props', () => {
    const periodicRefresh = {
      enabled: true,
      unit: TIME_UNITS.second,
      value: 2,
    };
    const wrapper = snapshotFactory({
      propsData: {
        periodicRefresh,
        label: 'Custom label',
        name: 'customName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
