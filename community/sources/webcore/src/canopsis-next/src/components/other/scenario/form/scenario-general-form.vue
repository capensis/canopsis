<template>
  <v-layout class="gap-3" column>
    <c-name-field
      v-field="form.name"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="$t('common.priority')">
        <c-priority-field v-field="form.priority" />
      </c-form-block-row>

      <c-form-block-row :label="$tc('common.trigger', 2)" required>
        <c-triggers-field
          v-field="form.triggers"
          with-additional-values
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.delay')">
        <c-duration-field
          v-field="form.delay"
          :label="$t('common.delay')"
          :units-label="$t('common.unit')"
          name="delay"
          clearable
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.disableDuringPeriods')">
        <c-disable-during-periods-field v-field="form.disable_during_periods" />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
