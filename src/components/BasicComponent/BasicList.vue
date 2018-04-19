<template>
    <ol>
        <li><slot name="header"></slot></li>
        <!-- vbind avec l'id de la base -->
        <li v-for="item in items">
          <list-brick :item="item">
              <div slot="reduced">
                <slot name="row"></slot>
              </div>
              <div slot="expanded">
                <slot name="expandedRow"></slot>
              </div>
          </list-brick>
        </li>
    </ol>
</template>

<script>
import ListBrick from './ListBrick.vue';

export default {
  name: 'BasicList',
  data() {
    return {
      isExpanded: false,
    };
  },
  components: { ListBrick },
  props: ['items'],
  methods: {
    onClick: () => {console.log('lol'); },
  },
};
</script>

<style scoped>

</style>
