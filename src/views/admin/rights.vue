<template lang="pug">
  v-container
    v-layout(v-for="(actions, key) in groupedActions", :key="key", row)
      v-flex(xs5)
        v-layout(row)
          v-flex(xs5)
            div {{ key }}
          v-flex(xs7)
            div(v-for="action in actions", :key="action._id") {{ action.desc }}
            div Action
            div Action
            div Action
      v-flex(xs7)
        v-layout(row)
          v-flex(xs12)
            div(v-for="action in actions", :key="`checkboxes-${action._id}`")
              v-checkbox
</template>

<script>
import entitiesActionMixin from '@/mixins/entities/action';
import entitiesRoleMixin from '@/mixins/entities/role';

export default {
  mixins: [entitiesActionMixin, entitiesRoleMixin],
  data() {
    return {
      pending: true,
    };
  },
  computed: {
    groupedActions() {
      return this.actions.reduce((acc, action) => {
        if (action.id.startsWith('view') || action.id.startsWith('userview')) {
          acc.view.push(action);
        } else if (action.id.startsWith('models')) {
          acc.technical.push(action);
        } else {
          acc.business.push(action);
        }

        return acc;
      }, { business: [], view: [], technical: [] });
    },
  },
  mounted() {
    this.fetchActionsList({
      params: { limit: 10000 },
    });

    this.fetchRolesList({
      params: { limit: 10000 },
    });
  },
};
</script>
