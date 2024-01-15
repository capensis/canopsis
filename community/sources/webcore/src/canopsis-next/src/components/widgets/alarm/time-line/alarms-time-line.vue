<template>
  <div class="timeline">
    <alarms-time-line-steps :steps="steps.data">
      <template #card="{ step }">
        <alarms-time-line-card
          :step="step"
          :is-html-enabled="isHtmlEnabled"
        />
      </template>
    </alarms-time-line-steps>
    <c-pagination
      :total="meta.total_count"
      :limit="meta.per_page"
      :page="meta.page"
      @input="updatePage"
    />
  </div>
</template>

<script>
import AlarmsTimeLineCard from './alarms-time-line-card.vue';
import AlarmsTimeLineSteps from './alarms-time-line-steps.vue';

export default {
  components: { AlarmsTimeLineSteps, AlarmsTimeLineCard },
  props: {
    steps: {
      type: Object,
      required: true,
    },
    isHtmlEnabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    meta() {
      return this.steps?.meta ?? {};
    },
  },
  methods: {
    updatePage(page) {
      this.$emit('update:page', page);
    },
  },
};
</script>

<style lang="scss" scoped>
.timeline {
  margin: 0 auto;
  width: 90%;
}
</style>
