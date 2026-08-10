<template>
  <v-row density="comfortable">
    <v-col cols="6">
      <c-autocomplete
        v-model="moduleName"
        v-model:search="searchInput"
        :items="modules"
        :loading="modulesPending"
        :placeholder="t('snmpRule.module')"
        class="pt-0"
        item-title="moduleName"
        item-value="moduleName"
        name="moduleName"
        hide-no-data
        hide-details
        required
        @update:model-value="selectModule"
      />
    </v-col>
    <v-col cols="6">
      <c-autocomplete
        v-model="mib"
        :items="moduleMibs"
        :loading="moduleMibsPending"
        :menu-props="{ offsetY: true }"
        class="pt-0"
        item-title="name"
        item-value="_id"
        name="oid.mibName"
        hide-no-data
        hide-details
        return-object
      />
    </v-col>
  </v-row>
</template>

<script>
import { find, sortBy, throttle } from 'lodash';
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import { MAX_LIMIT } from '@/constants';
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';
import { useSnmpMib } from '@/hooks/store/modules/snmp-mib';

export default {
  props: {
    form: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const modules = ref([]);
    const moduleMibs = ref([]);
    const searchInput = ref('');

    const { updateField } = useModelField(props, emit);
    const { fetchSnmpMibsList } = useSnmpMib();

    const moduleName = computed({
      get: () => props.form.moduleName,
      set: value => updateField('moduleName', value),
    });

    const mib = computed({
      get: () => props.form.mib,
      set: value => updateField('mib', value),
    });

    const { pending: modulesPending, handler: fetchModulesList } = usePendingHandler(async () => {
      const { data } = await fetchSnmpMibsList({
        params: {
          search: searchInput.value,

          limit: MAX_LIMIT,
          nodetype: 'notification',
          projection: 'moduleName',
          distinct: true,
        },
      });

      modules.value = sortBy(data, 'moduleName');
    });

    const throttledFetchModulesList = throttle(fetchModulesList, 500);

    const { pending: moduleMibsPending, handler: selectModule } = usePendingHandler(
      async (module) => {
        updateField('moduleName', module);

        const { data } = await fetchSnmpMibsList({
          params: {
            limit: MAX_LIMIT,
            nodetype: 'notification',
            moduleName: module,
          },
        });

        moduleMibs.value = sortBy(data, 'name');

        if (props.form.mib?.name) {
          updateField('mib', find(data, { name: props.form.mib.name }));
        }
      },
    );

    watch(searchInput, throttledFetchModulesList);

    onMounted(() => {
      fetchModulesList();

      if (props.form.moduleName) {
        selectModule(props.form.moduleName);
      }
    });

    onBeforeUnmount(() => throttledFetchModulesList.cancel());

    return {
      t,
      moduleName,
      mib,
      modules,
      modulesPending,
      moduleMibs,
      moduleMibsPending,
      searchInput,
      selectModule,
    };
  },
};
</script>
