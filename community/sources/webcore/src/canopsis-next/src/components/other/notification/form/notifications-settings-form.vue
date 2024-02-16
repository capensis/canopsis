<template>
  <v-layout column>
    <h4 class="text-subtitle-1">
      {{ $t('notificationSetting.instruction.header') }}
    </h4>
    <v-checkbox
      class="my-2"
      v-field="form.instruction.rate"
      :disabled="!form.instruction.rate_frequency.value"
      :label="$t('notificationSetting.instruction.rate')"
      color="primary"
      hide-details
    />
    <v-layout>
      <v-flex
        class="text-subtitle-1 mt-3"
        xs2
      >
        {{ $t('notificationSetting.instruction.rateFrequency') }}
      </v-flex>
      <v-flex xs3>
        <c-duration-field
          v-field="form.instruction.rate_frequency"
          :disabled="!form.instruction.rate"
          :required="form.instruction.rate"
          :label="$t('notificationSetting.instruction.duration')"
          long
        />
      </v-flex>
    </v-layout>
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
      required: true,
    },
  },
};
</script>
