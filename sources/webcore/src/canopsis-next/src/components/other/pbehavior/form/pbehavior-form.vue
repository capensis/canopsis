<template lang="pug">
  div
    pbehavior-general-form(:form="form.general", @input="updateField('general', $event)", :noFilter="noFilter")
    pbehavior-exdates-form.mt-2(
      v-show="form.general.rrule",
      :exdates="form.exdate",
      @input="updateField('exdate', $event)"
    )
    pbehavior-comments-form.mt-2(:comments="form.comments", @input="updateField('comments', $event)", :author="author")
</template>

<script>
import formMixin from '@/mixins/form';

import PbehaviorGeneralForm from './partials/pbehavior-general-form.vue';
import PbehaviorExdatesForm from './partials/pbehavior-exdates-form.vue';
import PbehaviorCommentsForm from './partials/pbehavior-comments-form.vue';

export default {
  components: {
    PbehaviorGeneralForm,
    PbehaviorExdatesForm,
    PbehaviorCommentsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    author: {
      type: String,
    },
    noFilter: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
