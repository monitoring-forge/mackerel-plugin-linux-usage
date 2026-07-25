package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/prometheus/procfs"
)

func buildProcfsTree(t *testing.T, procDir string) {
	testProcfsTree := []struct {
		path string
		body string
	}{
		{
			path: "stat",
			body: `cpu  18427395 4118474 10187089 6264581068 1094707 5493672 5806268 350653 0 0
cpu0 18427395 4118474 10187089 6264581068 1094707 5493672 5806268 350653 0 0
intr 4408091998 0 217 0 0 0 0 0 0 0 0 0 32 15 0 61685022 0 0 0 0 0 0 0 0 0 0 60024224 0 1482452849 87284685 0 99404096 2836623 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 16273431385
btime 1721797172
processes 9773755
procs_running 1
procs_blocked 0
softirq 6428064884 0 1343881500 732 1669197885 90021451 0 429383 0 6764 3324527169`,
		},
		{
			path: "loadavg",
			body: `0.13 0.06 0.02 2/107 1385894`,
		},
		{
			path: "1386385/net/snmp",
			body: `Ip: Forwarding DefaultTTL InReceives InHdrErrors InAddrErrors ForwDatagrams InUnknownProtos InDiscards InDelivers OutRequests OutDiscards OutNoRoutes ReasmTimeout ReasmReqds ReasmOKs ReasmFails FragOKs FragFails FragCreates
Ip: 2 64 219354148 0 0 0 0 0 163331259 87662751 0 20 202 6472 2713 202 36 0 72
Icmp: InMsgs InErrors InCsumErrors InDestUnreachs InTimeExcds InParmProbs InSrcQuenchs InRedirects InEchos InEchoReps InTimestamps InTimestampReps InAddrMasks InAddrMaskReps OutMsgs OutErrors OutDestUnreachs OutTimeExcds OutParmProbs OutSrcQuenchs OutRedirects OutEchos OutEchoReps OutTimestamps OutTimestampReps OutAddrMasks OutAddrMaskReps
Icmp: 5282410 1821 4 7567 184 0 0 101 5272265 2153 133 0 3 0 16892436 0 11617833 28 0 0 0 2177 5272398 0 0 0 0
IcmpMsg: InType0 InType3 InType5 InType8 InType11 InType13 InType17 OutType0 OutType3 OutType8 OutType11
IcmpMsg: 2153 7567 101 5272265 184 133 3 5272398 11617833 2177 28
Tcp: RtoAlgorithm RtoMin RtoMax MaxConn ActiveOpens PassiveOpens AttemptFails EstabResets CurrEstab InSegs OutSegs RetransSegs InErrs OutRsts InCsumErrors
Tcp: 1 200 120000 -1 128974 4576939 81904 414865 3 73695506 82233597 4175148 112 197074 0
Udp: InDatagrams NoPorts InErrors OutDatagrams RcvbufErrors SndbufErrors InCsumErrors IgnoredMulti MemErrors
Udp: 84353292 6 45 211916 45 0 0 0 0
UdpLite: InDatagrams NoPorts InErrors OutDatagrams RcvbufErrors SndbufErrors InCsumErrors IgnoredMulti MemErrors
UdpLite: 0 0 0 0 0 0 0 0 0`,
		},
		{
			path: "1386385/net/netstat",
			body: `TcpExt: SyncookiesSent SyncookiesRecv SyncookiesFailed EmbryonicRsts PruneCalled RcvPruned OfoPruned OutOfWindowIcmps LockDroppedIcmps ArpFilter TW TWRecycled TWKilled PAWSActive PAWSEstab DelayedACKs DelayedACKLocked DelayedACKLost ListenOverflows ListenDrops TCPHPHits TCPPureAcks TCPHPAcks TCPRenoRecovery TCPSackRecovery TCPSACKReneging TCPSACKReorder TCPRenoReorder TCPTSReorder TCPFullUndo TCPPartialUndo TCPDSACKUndo TCPLossUndo TCPLostRetransmit TCPRenoFailures TCPSackFailures TCPLossFailures TCPFastRetrans TCPSlowStartRetrans TCPTimeouts TCPLossProbes TCPLossProbeRecovery TCPRenoRecoveryFail TCPSackRecoveryFail TCPRcvCollapsed TCPBacklogCoalesce TCPDSACKOldSent TCPDSACKOfoSent TCPDSACKRecv TCPDSACKOfoRecv TCPAbortOnData TCPAbortOnClose TCPAbortOnMemory TCPAbortOnTimeout TCPAbortOnLinger TCPAbortFailed TCPMemoryPressures TCPMemoryPressuresChrono TCPSACKDiscard TCPDSACKIgnoredOld TCPDSACKIgnoredNoUndo TCPSpuriousRTOs TCPMD5NotFound TCPMD5Unexpected TCPMD5Failure TCPSackShifted TCPSackMerged TCPSackShiftFallback TCPBacklogDrop PFMemallocDrop TCPMinTTLDrop TCPDeferAcceptDrop IPReversePathFilter TCPTimeWaitOverflow TCPReqQFullDoCookies TCPReqQFullDrop TCPRetransFail TCPRcvCoalesce TCPOFOQueue TCPOFODrop TCPOFOMerge TCPChallengeACK TCPSYNChallenge TCPFastOpenActive TCPFastOpenActiveFail TCPFastOpenPassive TCPFastOpenPassiveFail TCPFastOpenListenOverflow TCPFastOpenCookieReqd TCPFastOpenBlackhole TCPSpuriousRtxHostQueues BusyPollRxPackets TCPAutoCorking TCPFromZeroWindowAdv TCPToZeroWindowAdv TCPWantZeroWindowAdv TCPSynRetrans TCPOrigDataSent TCPHystartTrainDetect TCPHystartTrainCwnd TCPHystartDelayDetect TCPHystartDelayCwnd TCPACKSkippedSynRecv TCPACKSkippedPAWS TCPACKSkippedSeq TCPACKSkippedFinWait2 TCPACKSkippedTimeWait TCPACKSkippedChallenge TCPWinProbe TCPKeepAlive TCPMTUPFail TCPMTUPSuccess TCPDelivered TCPDeliveredCE TCPAckCompressed TCPZeroWindowDrop TCPRcvQDrop TCPWqueueTooBig TcpTimeoutRehash
TcpExt: 0 0 0 81898 21 0 0 34 0 0 90519 0 0 0 7568 8535195 721 245083 0 376 26652719 12335141 19256128 0 44811 8 13598 1837 20 11 19 3870 46857 2000054 1 84 2614 80879 1193 1144572 889235 180586 0 8675 0 173249 249949 2530 171977 42 39405 11599 0 285298 0 0 0 0 10 49 110237 2 0 0 0 75914 45013 82657 0 0 0 0 7990532 0 0 0 0 10048028 269054 0 2699 1439 154 0 0 0 0 0 772 0 4 0 16539 1761 1761 4052 845523 47431623 12 1265 5 177 168 2177 3528 2 0 57 4 72 0 0 46890806 12838 31095 0 0 0 2440055
IpExt: InNoRoutes InTruncatedPkts InMcastPkts OutMcastPkts InBcastPkts OutBcastPkts InOctets OutOctets InMcastOctets OutMcastOctets InBcastOctets OutBcastOctets InCsumErrors InNoECTPkts InECT1Pkts InECT0Pkts InCEPkts ReasmOverlaps
IpExt: 0 0 0 0 130130 0 122453852496 34463927670 0 0 42629671 0 0 219430192 2626 192012 14571 0
MPTcpExt: MPCapableSYNRX MPCapableSYNTX MPCapableSYNACKRX MPCapableACKRX MPCapableFallbackACK MPCapableFallbackSYNACK MPFallbackTokenInit MPTCPRetrans MPJoinNoTokenFound MPJoinSynRx MPJoinSynAckRx MPJoinSynAckHMacFailure MPJoinAckRx MPJoinAckHMacFailure DSSNotMatching InfiniteMapRx OFOQueueTail OFOQueue OFOMerge NoDSSInWindow DuplicateData AddAddr EchoAdd PortAdd AddAddrDrop MPJoinPortSynRx MPJoinPortSynAckRx MPJoinPortAckRx MismatchPortSynRx MismatchPortAckRx RmAddr RmAddrDrop RmSubflow MPPrioTx MPPrioRx RcvPruned SubflowStale SubflowRecover
MPTcpExt: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0`,
		},
		{
			path: "1/stat",
			body: `1 (systemd) S 0 1 1 0 -1 4210944 1806767 1791322178 3043 1244804 72837 38109 8826700 1881909 20 0 1 0 0 243937280 1079 18446744073709551615 94308452241408 94308453570588 140734770805216 0 0 0 671173123 4096 1260 1 0 0 17 0 0 0 608 0 0 94308455669456 94308455923888 94308478599168 140734770810691 140734770810758 140734770810758 140734770810847 0`,
		},
		{
			path: "1386274/stat",
			body: `1386274 (top) R 1385833 1386274 1385833 34816 1386274 4210688 326 0 0 0 633 682 0 0 20 0 1 0 6316675639 53370880 1033 18446744073709551615 93978762792960 93978762902208 140728930957424 0 0 0 0 0 2147155711 0 0 0 17 0 0 0 0 0 0 93978765002960 93978765008728 93978778165248 140728930964940 140728930964944 140728930964944 140728930967531 0`,
		},
		{
			path: "1385833/status",
			body: `1385833 (bash) S 1385832 1385833 1385833 34816 1386321 4210944 15551 20717 8 30 14 4 710 766 20 0 1 0 6316535469 15155200 1346 18446744073709551615 94422119411712 94422120491960 140721051392896 0 0 0 65536 3670020 1266777851 0 0 0 17 0 0 0 2 0 0 94422122589456 94422122636676 94422130855936 140721051393543 140721051393549 140721051393549 140721051394030 0`,
		},
	}

	// Create a test procfs filesystem in a temporary directory
	err := os.MkdirAll(procDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test procfs directory: %v", err)
	}
	for _, file := range testProcfsTree {
		fullPath := filepath.Join(procDir, file.path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(file.body), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", fullPath, err)
		}
	}
	// Create symlink for /proc/self to /proc/1386385
	symlinkPath := filepath.Join(procDir, "self")
	if err := os.Symlink("1386385", symlinkPath); err != nil {
		t.Fatalf("Failed to create symlink for /proc/self: %v", err)
	}
}

// Test gaugeMetrics
func TestGaugeMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	procDir := filepath.Join(tmpDir, "proc")
	buildProcfsTree(t, procDir)

	// Set the workDir to the temporary directory
	plugin := LinuxUsagePlugin{
		workDir: tmpDir,
	}

	// Create a procfs.FS instance pointing to the temporary procfs directory
	pf, err := procfs.NewFS(procDir)
	if err != nil {
		t.Fatalf("Failed to create procfs FS: %v", err)
	}

	// Call gaugeMetrics and check for errors
	res, err := plugin.gaugeMetrics(pf)
	if err != nil {
		t.Fatalf("gaugeMetrics failed: %v", err)
	}

	// Check that the result contains expected keys
	expected := []struct {
		key   string
		value float64
	}{
		{"loadavg1", 0.13},
		{"loadavg5", 0.06},
		{"loadavg15", 0.02},
		{"all", 2},
		{"running", 1},
		{"active", 128974},
		{"passive", 4576939},
		{"overflows", 0},
		{"drops", 376},
	}
	for _, exp := range expected {
		if val, ok := res[exp.key]; !ok {
			t.Errorf("Expected key %s not found in result", exp.key)
		} else if val != exp.value {
			t.Errorf("Expected value for key %s: %v, got: %v", exp.key, exp.value, val)
		}
	}
}

// Test cpuMetrics
func TestCpuMetrics(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a procfs.FS instance pointing to the temporary procfs directory
	procDir := filepath.Join(tmpDir, "proc")
	err := os.MkdirAll(procDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test procfs directory: %v", err)
	}
	pf, err := procfs.NewFS(procDir)
	if err != nil {
		t.Fatalf("Failed to create procfs FS: %v", err)
	}

	// Create /proc/stat file with sample CPU stats
	statContent := `cpu  18428806 4118474 10188277 6264808253 1094741 5493853 5806514 350672 0 0
cpu0 18428806 4118474 10188277 6264808253 1094741 5493853 5806514 350672 0 0
intr 4408391516 0 217 0 0 0 0 0 0 0 0 0 32 15 0 61687274 0 0 0 0 0 0 0 0 0 0 60026721 0 1482503561 87307777 0 99407502 2836720 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 16273952183
btime 1721797172
processes 9774421
procs_running 1
procs_blocked 0
softirq 6428505879 0 1343933345 732 1669274374 90025047 0 430534 0 6764 3324835083`
	statPath := filepath.Join(procDir, "stat")
	if err := os.WriteFile(statPath, []byte(statContent), 0644); err != nil {
		t.Fatalf("Failed to write /proc/stat: %v", err)
	}

	plugin := LinuxUsagePlugin{
		workDir: tmpDir,
	}
	res, err := plugin.cpuMetrics(pf)
	if err != nil {
		t.Fatalf("cpuMetrics failed: %v", err)
	}
	// first time execution, res should be empty
	if len(res) != 0 {
		t.Errorf("Expected empty result on first execution, got: %v", res)
	}

	// Update the /proc/stat file to simulate a second execution
	statContent2 := `cpu  18429456 4118474 10188983 6264821901 1094743 5493871 5806556 350673 0 0
cpu0 18429456 4118474 10188983 6264821901 1094743 5493871 5806556 350673 0 0
intr 4408437410 0 217 0 0 0 0 0 0 0 0 0 32 15 0 61687420 0 0 0 0 0 0 0 0 0 0 60026881 0 1482509681 87321000 0 99407642 2836728 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 16274019462
btime 1721797172
processes 9774466
procs_running 2
procs_blocked 0
softirq 6428709196 0 1343939038 732 1669293799 90025278 0 431563 0 6764 3325012022`
	if err := os.WriteFile(statPath, []byte(statContent2), 0644); err != nil {
		t.Fatalf("Failed to write /proc/stat for second execution: %v", err)
	}

	res2, err := plugin.cpuMetrics(pf)
	if err != nil {
		t.Fatalf("cpuMetrics failed on second execution: %v", err)
	}
	// Check that the result contains expected keys
	expected := []struct {
		key string
		val float64
	}{
		{"user", 4.314063848238907},
		{"nice", 0.0},
		{"system", 4.685737041316406},
		{"idle", 90.5820667682313},
		{"iowait", 0.013274042610255611},
		{"irq", 0.11946638348988597},
		{"softirq", 0.27875489480812427},
		{"steal", 0.006637021305127806},
		{"guest", 0},
		{"guest_nice", 0},
	}
	for _, exp := range expected {
		if val, ok := res2[exp.key]; !ok {
			t.Errorf("Expected key %s not found in result on second execution", exp.key)
		} else if val != exp.val {
			t.Errorf("Expected value for key %s on second execution: %v, got: %v", exp.key, exp.val, val)
		}
	}
}
