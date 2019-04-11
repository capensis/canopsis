<template lang="pug">
  div.mt-1
    div(v-for="category in linkList", :key="category.cat_name")
      template(v-if="category.links.length")
        span.category.mr-2 {{ category.cat_name }}
        v-divider(light)
        div(v-for="(link, index) in category.links", :key="`links-${index}`")
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
      const links = this.category ?
        this.links.filter(({ cat_name: catName }) => catName === this.category) :
        this.links;

      return links.map((category) => {
        const categoryLinks = category.links.reduce((acc, link, index) => {
          if (typeof link === 'object' && link.link && link.label) {
            acc.push(link);
          } else {
            acc.push({
              label: `${category.cat_name} - ${index}`,
              link,
            });
          }

          return acc;
        }, []);

        return {
          cat_name: category.cat_name,
          links: categoryLinks,
        };
      });
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
