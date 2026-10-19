import { range } from 'lodash';

import { fakeObject } from '@unit/data/common';

export const fakeUsersForTreeview = ({
  fields = {
    _id: 'datatype.uuid',
    username: 'internet.userName',
    firstname: 'name.firstName',
    lastname: 'name.lastName',
    email: 'internet.email',
  },
  count = 1,
  fake = true,
  suffix = '',
  depths = 1,
  parentFields = {},
}) => range(count)
  .map((value) => {
    const localSuffix = `${suffix}${value + count * depths}`;
    const user = fakeObject({ fields, fake, suffix: localSuffix });

    return {
      ...user,
      ...parentFields,

      children:
        (value % 2 && depths)
          ? fakeUsersForTreeview({
            fields,
            count,
            fake,
            suffix: localSuffix,
            depths: depths - 1,
            parentFields: {
              parentKey: user._id,
            },
          })
          : [],
    };
  });
