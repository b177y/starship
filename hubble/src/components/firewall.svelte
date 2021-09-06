<script>
    export let rule
    export let index
    export let netgroups
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    function dispatchUpdate(){
        dispatch('updateRule', {
                index: index,
                rule: rule
            })
    }
    function groupChecked(e, group){
        console.log(group, e.target.checked)
        groups = rule.groups
        if (e.target.checked){ // add group
            groups.push(group)
        } else {
            for( var i = 0; i < groups.length; i++){ 
                if ( groups[i] === group) { 
                    groups.splice(i, 1); 
                }
            }
        }
        rule.groups = groups;
        dispatchUpdate()
    }
</script>

<main class="grid grid-cols-5 mt-2 flex items-start">
    <div class="flex flex-row">
        <label for="port" class="my-auto">Port</label>
        <input id="port" name="port" type="text" bind:value={rule.port}
            class="w-1/2 h-7 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto mx-4 bg-gray-100"/>
    </div>
    <div class="flex flex-row">
        <label for="protocol" class="my-auto">Protocol</label>
        <select bind:value={rule.proto} on:blur={dispatchUpdate}
            class="w-1/2 h-7 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto mx-4">
            <option value="any">Any</option>
            <option value="icmp">ICMP</option>
            <option value="udp">UDP</option>
            <option value="tcp">TCP</option>
        </select>
    </div>
    <div class="flex flex-row">
        <label for="any" class="my-auto pr-4">Any Host</label>
        <input type="checkbox" bind:checked={rule.any} class="my-auto" />
    </div>
    {#if !rule.any}
    <div id="groups">
        <h1><b>Groups</b></h1>
        {#each netgroups as group}
            <div class="flex flex-row">
                <input type="checkbox" id={group} class="my-auto"
                    checked={rule.groups.includes(group)} on:change={(e) => groupChecked(e, group)}/>
                <label for="{group}" class="my-auto pl-2">{group}</label>
            </div>
        {/each}
    </div>
    {:else}
        <div class="w-1/5"></div>
    {/if}
    <div class="flex justify-end">
    <button on:click={() => dispatch('deleteRule')}
        class="bg-red-400 text-white rounded-xl p-1 h-7 w-1/2">Delete</button>
    </div>
</main>
