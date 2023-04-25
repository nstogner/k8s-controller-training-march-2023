# Questions

- What component handles the lifecycle of a Node object?
    - Kubelet can create Node objects. Or they can be created out of band.
    - The Node Controller is responsible to deleting Node objects for VMs that no longer exist in the cloud provider.
- Are objects persisted and then broadcasted or visa versa?
    - Etcd's built in watch mechanisms are used to serve watches.
    - So: Persistence, then watch broadcast.

## Future Followups

- Can you filter events based on type on the Server-side?

- If a client requests a watch and the connection is severed is there a way it can re-attach to ensure it has missed nothing?

- What are these update errors that I commonly see (referring to version conflicts)?

- Dream controller example.

- Cloud API Controllers

- Show how to bypass the cache

- Show the code that does the `.Owns()` filtering.

- Provide some examples of diff computing / hashing / patching to ensure desired state.