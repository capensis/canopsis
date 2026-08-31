<template>
  <div>
    <c-layout>
      <c-pbehavior-type-field
        v-model="activeOnPbhField"
        :label="t('remediation.pattern.tabs.pbehaviorTypes.fields.activeOnTypes')"
        :is-item-disabled="isActiveItemDisabled"
        chips
        multiple
      />
    </c-layout>
    <c-layout>
      <c-pbehavior-type-field
        v-model="disabledOnPbhField"
        :label="t('remediation.pattern.tabs.pbehaviorTypes.fields.disabledOnTypes')"
        :is-item-disabled="isDisabledItemDisabled"
        chips
        multiple
      />
    </c-layout>
  </div>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ['update:form'],
  setup(props, { emit }) {
    const { t } = useI18n();
    const createFieldModel = field => computed({
        get: () => props.form[field],
        set: value => emit('update:form', {
            ...props.form,
            [field]: value,
          }),
      });
    const activeOnPbhField = createFieldModel('active_on_pbh');
    const disabledOnPbhField = createFieldModel('disabled_on_pbh');

    const isActiveItemDisabled = item => props.form.disabled_on_pbh.includes(item._id);
    const isDisabledItemDisabled = item => props.form.active_on_pbh.includes(item._id);

    return {
      t,
      activeOnPbhField,
      disabledOnPbhField,
      isActiveItemDisabled,
      isDisabledItemDisabled,
    };
  },
};
</script>
