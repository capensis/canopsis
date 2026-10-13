import Handlebars from 'handlebars';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars/registers';

describe('Handlebars helper registers', () => {
  const helperName = 'test-helper';
  const helper = jest.fn(() => 'result');

  let instance;

  beforeEach(() => {
    instance = Handlebars.create();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('Helper stays registered until the last owner unregisters it', () => {
    const unregisterHelperSpy = jest.spyOn(instance, 'unregisterHelper');

    registerHelper(helperName, helper, instance);
    registerHelper(helperName, helper, instance);

    unregisterHelper(helperName, instance);

    expect(instance.helpers[helperName]).toBe(helper);
    expect(unregisterHelperSpy).not.toHaveBeenCalled();

    unregisterHelper(helperName, instance);

    expect(instance.helpers[helperName]).toBeUndefined();
    expect(unregisterHelperSpy).toHaveBeenCalledTimes(1);
    expect(unregisterHelperSpy).toHaveBeenCalledWith(helperName);
  });

  test('Helper owners are counted separately for each Handlebars instance', () => {
    const anotherInstance = Handlebars.create();

    registerHelper(helperName, helper, instance);
    registerHelper(helperName, helper, anotherInstance);

    unregisterHelper(helperName, instance);

    expect(instance.helpers[helperName]).toBeUndefined();
    expect(anotherInstance.helpers[helperName]).toBe(helper);

    unregisterHelper(helperName, anotherInstance);

    expect(anotherInstance.helpers[helperName]).toBeUndefined();
  });

  test('Helper registered by external code is not removed', () => {
    const externalHelper = jest.fn(() => 'external result');

    instance.registerHelper(helperName, externalHelper);

    const unregisterHelperSpy = jest.spyOn(instance, 'unregisterHelper');

    registerHelper(helperName, helper, instance);
    unregisterHelper(helperName, instance);

    expect(instance.helpers[helperName]).toBe(externalHelper);
    expect(unregisterHelperSpy).not.toHaveBeenCalled();
  });
});
