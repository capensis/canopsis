<template>
  <div>
    <snmp-rule-form-field-title :label="$t('snmpRule.oid')" />
    <v-layout wrap>
      <v-flex
        class="pr-1"
        xs6
      >
        <v-autocomplete
          v-validate="'required'"
          :value="form.moduleName"
          :items="modules"
          :search-input.sync="searchInput"
          :loading="modulesPending"
          :placeholder="$t('snmpRule.module')"
          :error-messages="errors.collect('moduleName')"
          class="pt-0"
          item-text="moduleName"
          item-value="_id"
          name="moduleName"
          autofocus
          hide-no-data
          hide-details
          @change="selectModule"
        />
      </v-flex>
      <v-flex
        class="pl-1"
        xs6
      >
        <v-autocomplete
          :value="form.mib"
          :items="moduleMibs"
          :loading="moduleMibsPending"
          :menu-props="{ offsetY: true }"
          :error-messages="errors.collect(mibName)"
          :name="mibName"
          class="pt-0"
          item-text="name"
          item-value="_id"
          hide-no-data
          hide-details
          return-object
          @change="selectMib"
        />
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import { find, sortBy, throttle } from 'lodash';
import {
  ref,
  watch,
  nextTick,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { MAX_LIMIT } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { usePendingHandler } from '@/hooks/query/pending';
import { useSnmpMib } from '@/hooks/store/modules/snmp-mib';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';

import SnmpRuleFormFieldTitle from './snmp-rule-form-field-title.vue';

export default {
  inject: ['$validator'],
  components: { SnmpRuleFormFieldTitle },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
    const mibName = 'oid.mibName';

    const modules = ref([]);
    const moduleMibs = ref([]);
    const searchInput = ref('');

    const { updateField } = useModelField(props, emit);
    const { fetchSnmpMibList } = useSnmpMib();

    const {
      attachRequiredRule,
      detachRequiredRule,
      validateRequiredRule,
    } = useValidationAttachRequired(mibName);

    const {
      pending: modulesPending,
      handler: fetchModulesList,
    } = usePendingHandler(async () => {
      const { data } = await fetchSnmpMibList({
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

    const {
      pending: moduleMibsPending,
      handler: selectModule,
    } = usePendingHandler(async (module) => {
      updateField('moduleName', module);

      const { data } = await fetchSnmpMibList({
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
    });

    /**
     * Select a MIB and update the form field.
     *
     * @param {Object} mib - The MIB object to be selected.
     * @returns {void}
     */
    const selectMib = (mib) => {
      updateField('mib', mib);

      nextTick(validateRequiredRule);
    };

    /**
     * Check if a MIB name exists.
     *
     * @returns {boolean} Returns true if the MIB name exists, otherwise false.
     */
    const isMibNameExists = () => !!props.form.mib?.name;

    watch(searchInput, throttledFetchModulesList);

    onMounted(() => {
      if (props.form.moduleName) {
        selectModule(props.form.moduleName);
      }

      attachRequiredRule(isMibNameExists);
    });

    onBeforeUnmount(detachRequiredRule);

    return {
      mibName,

      modules,
      modulesPending,
      moduleMibs,
      moduleMibsPending,
      searchInput,

      selectModule,
      selectMib,
    };
  },
};
</script>
