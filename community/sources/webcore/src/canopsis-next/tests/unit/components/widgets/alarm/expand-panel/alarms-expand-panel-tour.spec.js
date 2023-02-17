import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import { DEFAULT_LOCALE, THEMES_NAMES } from '@/config';
import { GROUPS_NAVIGATION_TYPES, TOURS } from '@/constants';

import AlarmsExpandPanelTour from '@/components/widgets/alarm/expand-panel/alarms-expand-panel-tour.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(AlarmsExpandPanelTour, {
  localVue,

  ...options,
});

const tours = {
  [Faker.datatype.string()]: Faker.datatype.boolean(),
  [TOURS.alarmsExpandPanel]: false,
};
const currentUser = {
  _id: Faker.datatype.string(),
  defaultview: '',
  email: Faker.internet.email(),
  enable: Faker.datatype.boolean(),
  firstname: Faker.name.firstName(),
  lastname: Faker.name.lastName(),
  name: Faker.name.firstName(),
  password: Faker.datatype.string(),
  role: '',
  ui_groups_navigation_type: GROUPS_NAVIGATION_TYPES.sideBar,
  ui_language: DEFAULT_LOCALE,
  ui_tours: tours,
  ui_theme: THEMES_NAMES.canopsis,
};

const fetchCurrentUser = jest.fn();
const authModule = {
  name: 'auth',
  getters: {
    currentUser,
  },
  actions: {
    fetchCurrentUser,
  },
};
const updateCurrentUser = jest.fn();
const userModule = {
  name: 'user',
  actions: {
    updateCurrentUser,
  },
};

const store = createMockedStoreModules([
  userModule,
  authModule,
]);

const selectSkipButton = wrapper => wrapper.find('.v-step__button-skip');
const selectPreviousButton = wrapper => wrapper.find('.v-step__button-previous');
const selectNextButton = wrapper => wrapper.find('.v-step__button-next');
const selectStopButton = wrapper => wrapper.find('.v-step__button-stop');

describe('alarms-expand-panel-tour', () => {
  jest.useFakeTimers();

  afterAll(() => {
    jest.useRealTimers();
  });

  afterEach(() => {
    updateCurrentUser.mockReset();
    fetchCurrentUser.mockReset();
  });

  it('Next callback called after trigger next button', async () => {
    const onNextStep = jest.fn();
    const wrapper = snapshotFactory({
      propsData: {
        callbacks: {
          onNextStep,
        },
      },
    });

    jest.runAllTimers();
    await flushPromises();
    const nextButton = selectNextButton(wrapper);
    nextButton.trigger('click');

    jest.runAllTimers();
    await flushPromises();

    expect(onNextStep).toHaveBeenCalledTimes(1);
  });

  it('Previous callback called after trigger previous button', async () => {
    const onPreviousStep = jest.fn();
    const wrapper = snapshotFactory({
      propsData: {
        callbacks: {
          onPreviousStep,
        },
      },
    });

    jest.runAllTimers();
    await flushPromises();

    const nextButton = selectNextButton(wrapper);
    nextButton.trigger('click');
    jest.runAllTimers();
    await flushPromises();

    const previousButton = selectPreviousButton(wrapper);
    previousButton.trigger('click');
    jest.runAllTimers();
    await flushPromises();

    expect(onPreviousStep).toHaveBeenCalledTimes(1);
  });

  it('Skip callback called after trigger skip button', async () => {
    const onSkip = jest.fn();
    const wrapper = snapshotFactory({
      store,
      propsData: {
        callbacks: {
          onSkip,
        },
      },
    });

    jest.runAllTimers();
    await flushPromises();

    const skipButton = selectSkipButton(wrapper);
    skipButton.trigger('click');
    jest.runAllTimers();
    await flushPromises();

    expect(onSkip).toHaveBeenCalledTimes(1);
  });

  it('Finish callback called after finish tour', async () => {
    const onFinish = jest.fn();
    const wrapper = snapshotFactory({
      store,
      propsData: {
        callbacks: {
          onFinish,
        },
      },
    });

    jest.runAllTimers();
    await flushPromises();
    selectNextButton(wrapper).trigger('click');

    jest.runAllTimers();
    await flushPromises();
    selectNextButton(wrapper).trigger('click');

    jest.runAllTimers();
    await flushPromises();

    selectStopButton(wrapper).trigger('click');
    jest.runAllTimers();
    await flushPromises();

    expect(onFinish).toHaveBeenCalledTimes(1);
  });

  it('User updated after skip', async () => {
    const wrapper = snapshotFactory({ store });

    jest.runAllTimers();
    await flushPromises();
    const skipButton = selectSkipButton(wrapper);
    skipButton.trigger('click');

    jest.runAllTimers();
    await flushPromises();

    expect(updateCurrentUser).toHaveBeenCalledTimes(1);
    expect(updateCurrentUser).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          ...currentUser,
          ui_tours: {
            ...tours,
            [TOURS.alarmsExpandPanel]: true,
          },
        },
      },
      undefined,
    );
    expect(fetchCurrentUser).toHaveBeenCalledTimes(1);
  });

  it('User updated after finish tour', async () => {
    const wrapper = snapshotFactory({ store });

    jest.runAllTimers();
    await flushPromises();
    selectNextButton(wrapper).trigger('click');

    jest.runAllTimers();
    await flushPromises();
    selectNextButton(wrapper).trigger('click');

    jest.runAllTimers();
    await flushPromises();

    selectStopButton(wrapper).trigger('click');
    jest.runAllTimers();
    await flushPromises();

    expect(updateCurrentUser).toHaveBeenCalledTimes(1);
    expect(updateCurrentUser).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          ...currentUser,
          ui_tours: {
            ...tours,
            [TOURS.alarmsExpandPanel]: true,
          },
        },
      },
      undefined,
    );
    expect(fetchCurrentUser).toHaveBeenCalledTimes(1);
  });

  it('Renders `alarms-expand-panel-tour` with first step', async () => {
    const wrapper = snapshotFactory();

    jest.runAllTimers();
    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel-tour` with second step', async () => {
    const wrapper = snapshotFactory();

    jest.runAllTimers();
    await flushPromises();
    const nextStepButton = selectNextButton(wrapper);

    nextStepButton.trigger('click');
    jest.runAllTimers();
    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel-tour` with third step', async () => {
    const wrapper = snapshotFactory();

    jest.runAllTimers();
    await flushPromises();
    selectNextButton(wrapper).trigger('click');

    jest.runAllTimers();
    await flushPromises();
    selectNextButton(wrapper).trigger('click');

    jest.runAllTimers();
    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel-tour` after finish', async () => {
    const wrapper = snapshotFactory({ store });

    jest.runAllTimers();
    await flushPromises();
    selectNextButton(wrapper).trigger('click');

    jest.runAllTimers();
    await flushPromises();
    selectNextButton(wrapper).trigger('click');

    jest.runAllTimers();
    await flushPromises();

    const stopButton = selectStopButton(wrapper);
    stopButton.trigger('click');
    jest.runAllTimers();
    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
