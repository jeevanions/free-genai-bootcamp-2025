import React from 'react';
import Select from 'react-select';
import type { Group } from '../types/api';

interface CategorySelectProps {
  groups: Group[];
  isLoading: boolean;
  value: Group | null;
  onChange: (group: Group | null) => void;
  onCreateGroup: (name: string) => Promise<void>;
}

export const CategorySelect: React.FC<CategorySelectProps> = ({
  groups,
  isLoading,
  value,
  onChange,
  onCreateGroup,
}) => {
  const options = groups.map(group => ({
    value: group.id,
    label: `${group.name} (${group.word_count} words)`,
    group,
  }));

  return (
    <Select
      className="w-full"
      options={options}
      isLoading={isLoading}
      value={value ? {
        value: value.id,
        label: `${value.name} (${value.word_count} words)`,
        group: value,
      } : null}
      onChange={(option) => onChange(option ? option.group : null)}
      onCreateOption={onCreateGroup}
      isClearable
      placeholder="Select or create a category..."
      formatCreateLabel={(inputValue) => `Create category "${inputValue}"`}
    />
  );
};