package ffiwrapper

import (
	"bufio"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	"io"
	"io/ioutil"
	"os"
)

var PledgeSectorPath string
var PledgeCidPath string

func (sb *Sealer) isPledgeRequest(entries int, ps abi.UnpaddedPieceSize) bool {
	return entries == 0 && ps == 34091302912
}

func (sb *Sealer) initEnvForPledge() {
	CCPath, ok := os.LookupEnv("CC_SECTOR_PATH")
	if !ok {
		panic("CC_SECTOR_PATH not found")
	}
	PledgeSectorPath = CCPath + "/sector"
	PledgeCidPath = CCPath + "/cid"
}

func (sb *Sealer) pledgeSectorExists() bool {
	_, err := os.Stat(PledgeSectorPath)
	if err != nil {
		log.Errorf("pledge sector does not exist: %v", err)
		return false
	}
	_, err = os.Stat(PledgeCidPath)
	if err != nil {
		log.Errorf("pledge cid does not exist: %v", err)
		return false
	}
	return true
}

func (sb *Sealer) getPledgeCid() (cid.Cid, error) {
	data, err := ioutil.ReadFile(PledgeCidPath)
	if err != nil {
		return cid.Undef, err
	}
	return cid.Parse(data)
}

func (sb *Sealer) setPledgeCid(pieceCid cid.Cid) error {
	data, err := pieceCid.MarshalBinary()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(PledgeCidPath, data, 0644)
}

func (sb *Sealer) fromPledgeSector(unsealed string) error {
	pledgeFile, err := os.Open(PledgeSectorPath)
	if err != nil {
		return err
	}
	defer pledgeFile.Close()
	reader := bufio.NewReader(pledgeFile)
	unsealedFile, err := os.OpenFile(unsealed, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer unsealedFile.Close()
	writer := bufio.NewWriter(unsealedFile)
	_, err = io.Copy(writer, reader)
	return err
}

func (sb *Sealer) toPledgeSector(unsealed string) error {
	unsealedFile, err := os.Open(unsealed)
	if err != nil {
		return err
	}
	defer unsealedFile.Close()
	reader := bufio.NewReader(unsealedFile)
	pledgeFile, err := os.OpenFile(PledgeSectorPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer pledgeFile.Close()
	writer := bufio.NewWriter(pledgeFile)
	_, err = io.Copy(writer, reader)
	return err
}

func (sb *Sealer) cleanPledgeSector() (err error) {
	err = os.Remove(PledgeSectorPath)
	err = os.Remove(PledgeCidPath)
	return err
}
