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
      <c-form-block-row :label="$t('modals.createPbehaviorReason.fields.description')">
        <c-description-field
          v-field="form.description"
          name="description"
          required
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
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
