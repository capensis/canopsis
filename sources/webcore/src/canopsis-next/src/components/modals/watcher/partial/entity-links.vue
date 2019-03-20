<template lang="pug">
  div.mt-1
    span.category.mr-2 {{ category || $t('common.links') }}
    v-divider(light)
    div(v-for="(link, index) in linkList", :key="`links-${index}`")
      div.pa-2.text-xs-right
        a(:href="link.link", target="_blank") {{ link.label }}
</template>

<script>
export default {
  props: {
    links: {
      type: Array,
      default: () => [],
    },
    category: {
      type: String,
      default: null,
    },
  },
  computed: {
    /*
    * The linkbuilders used to return the links directly as
    * strings. They can now also return objects with the
    * properties 'label' and 'link', allowing to change the link's
    * label.
    * The following code converts the "legacy" representation
    * (strings) into the "new" representation, so they can be
    * displayed in the same manner by the template.
    */
    linkList() {
      const filteredLinks = this.category ?
        this.links.filter(({ cat_name: catName }) => catName === this.category) :
        this.links;

      return filteredLinks[0].links.reduce((acc, link, index) => {
        if (typeof link === 'object' && link.link && link.label) {
          acc.push(link);
        } else {
          acc.push({
            label: `${this.category} - ${index}`,
            link,
          });
        }

        return acc;
      }, []);
    },
  },
};
</script>

<style lang="scss" scoped>
  .category {
    display: inline-block;

    &:first-letter {
      text-transform: uppercase;
    }
  }
</style>
