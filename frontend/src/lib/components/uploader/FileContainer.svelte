<script lang="ts">
import { truncate } from "$lib/text";
import type { Uploadable } from "$lib/types";
import Divider from "./Divider.svelte";
import * as Card from "$lib/components/ui/card";

export let file: Uploadable | undefined;

const fullName = (file && file.data) ? file.data.name : 'loading';
const name = truncate(fullName, 40);
</script>

{#if file === undefined}
  <div>loading</div>
{:else}
  <Card.Root>
    <Card.Header>
      <div class="flex flex-row justify-between">
        <div class="basis-3/4 text-xl font-bold">
          {name}
        </div>
        <div class="basis-1/4 text-right">
          <slot name="buttons"/>
        </div>
      </div>
    </Card.Header>
    <Card.Content>
      <slot/>
    </Card.Content>
  </Card.Root>
{/if}
