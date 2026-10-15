import { generateShallowRenderer } from '@unit/utils/vue';

import ScenarioActionCardCompiledTemplate from '@/components/other/scenario/partials/scenario-action-card-compiled-template.vue';

describe('scenario-action-card-compiled-template', () => {
  const factory = generateShallowRenderer(ScenarioActionCardCompiledTemplate);

  it('renders Go template variables as literal text', () => {
    const template = 'Trigger: {{ .AdditionalData.Trigger }}\n<strong>plain text</strong>';
    const wrapper = factory({ propsData: { template } });

    expect(wrapper.text()).toBe(template);
    expect(wrapper.find('strong').exists()).toBe(false);
  });
});
