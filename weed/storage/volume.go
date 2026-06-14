package storage

import (
	"sync"
)

// writeNeedle writes a needle to the volume and updates the index.
func (v *Volume) writeNeedle(n *Needle) (uint32, error) {
	v.accessLock.Lock()
	defer v.accessLock.Unlock()

	// Perform the write to the .dat file
	offset, size, err := v.appendNeedle(n)
	if err != nil {
		return 0, err
	}

	// Update the needle map index synchronously to ensure read-after-write consistency
	err = v.nm.Put(n.Id, offset, n.Size)
	return size, err
}

// deleteNeedle marks a needle as deleted in the volume and updates the index.
func (v *Volume) deleteNeedle(n *Needle) (uint32, error) {
	v.accessLock.Lock()
	defer v.accessLock.Unlock()

	// Append tombstone to .dat file
	offset, size, err := v.appendNeedle(n)
	if err != nil {
		return 0, err
	}

	// Update the needle map index to reflect deletion (tombstone)
	err = v.nm.Delete(n.Id)
	return size, err
}