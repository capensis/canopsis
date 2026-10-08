<template>
  <v-layout
    class="gap-3"
    column
  >
    <c-enabled-field
      v-field="form.visible"
      :label="$t('pbehavior.visible')"
      with-background
    />
    <c-name-field
      v-field="form.name"
      autofocus
      required
    />
    <c-form-block>
      <c-form-block-row :label="$t('modals.createPbehaviorException.fields.description')">
        <c-description-field
          v-field="form.description"
          name="description"
          required
        />
      </c-form-block-row>

      <c-form-block-row
        :label="$t('pbehavior.exdates.title')"
        :align-center="!form.exdates.length"
        indented
      >
        <pbehavior-exceptions-field
          v-field="form.exdates"
          with-exdate-type
          required
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import PbehaviorExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exceptions-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorExceptionsField },
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
