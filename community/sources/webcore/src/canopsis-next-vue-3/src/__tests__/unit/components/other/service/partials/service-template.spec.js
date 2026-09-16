import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import ServiceTemplate from '@/components/other/service/partials/service-template.vue';

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
  'c-compiled-template': CCompiledTemplate,
  'service-entities-list': true,
};

const selectEntitiesList = wrapper => wrapper.find('service-entities-list-stub');

describe('service-template', () => {
  const service = {
    _id: 'service-id',
  };
  const modalTemplate = `
    Id: {{ entity._id }}
    Entities: {{ entities name="entity._id" }}
    <br />
  `;

  const { authModule } = createAuthModule();

  const store = createMockedStoreModules([authModule]);

  const factory = generateRenderer(ServiceTemplate, {
    store,

    stubs,
    propsData: {
      service,
      options: {},
    },
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Refresh applied after triggers entities list', async () => {
    const wrapper = factory({
      propsData: {
        serviceEntities: [{}, {}],
        widgetParameters: {
          modalTemplate,
        },
      },
    });

    await flushPromises();

    const entitiesList = selectEntitiesList(wrapper);

    await entitiesList.triggerCustomEvent('refresh');

    expect(wrapper).toHaveBeenEmit('refresh');
  });

  test('Add action applied after triggers entities list', async () => {
    const wrapper = factory({
      propsData: {
        serviceEntities: [{}, {}],
        widgetParameters: {
          modalTemplate,
        },
      },
    });

    await flushPromises();

    await selectEntitiesList(wrapper).triggerCustomEvent('add:action');

    expect(wrapper).toHaveBeenEmit('add:action');
  });

  test('Pagination updated after triggers entities list', async () => {
    const wrapper = factory({
      propsData: {
        serviceEntities: [{}, {}],
        widgetParameters: {
          modalTemplate,
        },
      },
    });

    await flushPromises();

    const entitiesList = selectEntitiesList(wrapper);

    const newPagination = {
      page: Faker.datatype.number(),
    };

    await entitiesList.triggerCustomEvent('update:options', newPagination);

    expect(wrapper).toEmit('update:options', newPagination);
  });
});
