import Handlebars from 'handlebars';
import { createPinia, setActivePinia } from 'pinia';

import { useAuthStore } from '@/stores/modules/auth';
import { useTemplateVarsStore } from '@/stores/modules/template-vars';

import { compile } from '@/helpers/handlebars/compilers';
import { registerTemplate, runTemplate } from '@/helpers/handlebars/templates';

vi.unmock('@/stores/modules/template-vars');

describe.each(['compile', 'runTemplate'])('%s context', (method) => {
  let render;

  beforeEach(() => {
    setActivePinia(createPinia());

    const instance = Handlebars.create();

    render = (template, context) => {
      if (method === 'compile') {
        return compile(template, context, instance);
      }

      registerTemplate('test', template, instance);

      return runTemplate('test', context, instance);
    };
  });

  test('renders configured variables alongside user and explicit context', async () => {
    useTemplateVarsStore().items = { Location: 'Paris', System: { ENV_TEST: 'production' } };
    vi.spyOn(useAuthStore(), 'currentUser', 'get').mockReturnValue({ name: 'Alice' });

    const result = await render(
      '{{env.Location}} / {{env.System.ENV_TEST}} / {{user.name}} / {{version}}',
      { version: '26.10' },
    );

    expect(result).toBe('Paris / production / Alice / 26.10');
  });

  test('renders with no configured variables', async () => {
    expect(await render('{{env.Location}}Version {{version}}', { version: '26.10' })).toBe(
      'Version 26.10',
    );
  });

  test('allows explicit context to override configured variables', async () => {
    useTemplateVarsStore().items = { Location: 'Paris' };

    expect(await render('{{env.Location}}', { env: { Location: 'London' } })).toBe('London');
  });
});
