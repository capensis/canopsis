import Handlebars from 'handlebars';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars/registers';

describe('Handlebars helper registrations', () => {
  test('Keeps a shared helper registered until its last owner unregisters', () => {
    const instance = Handlebars.create();
    const helper = vi.fn(() => 'value');

    registerHelper('shared', helper, instance);
    registerHelper('shared', helper, instance);

    unregisterHelper('shared', instance);

    expect(instance.helpers.shared).toBe(helper);

    unregisterHelper('shared', instance);

    expect(instance.helpers.shared).toBeUndefined();
  });

  test('Does not remove a helper owned by an external registry', () => {
    const instance = Handlebars.create();
    const helper = vi.fn(() => 'value');

    instance.registerHelper('external', helper);
    registerHelper('external', helper, instance);
    unregisterHelper('external', instance);

    expect(instance.helpers.external).toBe(helper);
  });
});
